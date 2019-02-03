package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type Contest struct {
	// only c and id is populated before load()
	c                *Core
	id               primitive.ObjectID
	lock             sync.RWMutex
	running          bool
	loaded           bool
	schedules        []*contestSchedule
	scheduleTimer    *time.Timer
	state            string
	updateChan       chan mongo.WriteModel
	playerUpdateChan chan mongo.WriteModel
	context          context.Context
	cancelFunc       func()

	players map[primitive.ObjectID]*ContestPlayer

	// ranklist is NOT covcered by lock; it has its own synchronization mechanism
	ranklist ContestRanklist
	rankcomp ContestRankComp
}

type contestSchedule struct {
	t time.Time
	f func()
}

func newState() string {
	var state [16]byte
	rand.Read(state[:])
	return hex.EncodeToString(state[:])
}

// Call exactly once when the contest gets loaded into memory.
func (c *Contest) load(contestModel *model.Contest) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.WithField("contestId", c.id).Info("Loading contest")
	c.running = contestModel.Running
	c.state = contestModel.State
	c.loaded = true
	c.context, c.cancelFunc = context.WithCancel(context.Background())
	c.updateChan = make(chan mongo.WriteModel, 100)
	c.playerUpdateChan = make(chan mongo.WriteModel, 1000)
	switch contestModel.RanklistType {
	case "realtime":
		c.ranklist = &ContestRealTimeRanklist{c: c}
	default:
		c.ranklist = ContestDummyRanklist{}
	}
	c.ranklist.Load()
	switch contestModel.RanklistComp {
	case "maxsum":
		c.rankcomp = ContestRankCompMaxScoreSum{}
	case "lastsum":
		c.rankcomp = ContestRankCompLastSum{}
	case "acm":
		c.rankcomp = ContestRankCompACM{}
	default:
		c.rankcomp = ContestDummyRankComp{}
	}
	// Bring up the writer
	go handleWrites(c.context, c.updateChan, c.c.mongodb.Collection("problemset"), c.id)
	go handleWrites(c.context, c.playerUpdateChan, c.c.mongodb.Collection("contest_player"), c.id)

	// Load all schedules
	for id, scheduleModel := range contestModel.Schedule {
		if scheduleModel.Done {
			continue
		}
		var f func()
		switch scheduleModel.Type {
		case "start":
			f = func(c *Contest, id int) func() {
				return func() {
					c.running = true
					model := mongo.NewUpdateOneModel()
					model.SetFilter(bson.D{{"_id", c.id}, {"contest.state", c.state}})
					c.state = newState()
					model.SetUpdate(bson.D{{"$set", bson.D{{fmt.Sprintf("contest.schedule.%d.done", id), true}, {"contest.running", true}, {"contest.state", c.state}}}})
					c.updateChan <- model
					log.WithField("contestId", c.id).Debug("Contest started")
				}
			}(c, id)
		case "stop":
			f = func(c *Contest, id int) func() {
				return func() {
					c.running = false
					model := mongo.NewUpdateOneModel()
					model.SetFilter(bson.D{{"_id", c.id}, {"contest.state", c.state}})
					c.state = newState()
					model.SetUpdate(bson.D{{"$set", bson.D{{fmt.Sprintf("contest.schedule.%d.done", id), true}, {"contest.running", false}, {"contest.state", c.state}}}})
					c.updateChan <- model
					log.WithField("contestId", c.id).Debug("Contest stopped")
				}
			}(c, id)
		}
		c.schedules = append(c.schedules, &contestSchedule{
			t: scheduleModel.StartTime,
			f: f,
		})
	}
	c.sortSchedules()

	// Load all players
	c.players = make(map[primitive.ObjectID]*ContestPlayer)
	var (
		cursor mongo.Cursor
		err    error
	)
	if cursor, err = c.c.mongodb.Collection("contest_player").Find(c.context, bson.D{{"contest", c.id}}); err != nil {
		log.WithField("contestId", c.id).Warning("Failed to load contest players: ", err)
	}
	for cursor.Next(c.context) {
		var contestPlayerModel *model.ContestPlayer
		if err = cursor.Decode(&contestPlayerModel); err != nil {
			log.WithField("contestId", c.id).Warning("Failed to load a contest player: ", err)
		} else {
			c.loadPlayer(contestPlayerModel)
		}
	}
	log.WithField("contestId", c.id).WithField("playerCount", len(c.players)).Info("Loaded players")
	go c.startSchedule()
}

type scheduleSorter struct {
	*Contest
}

func (s scheduleSorter) Len() int {
	return len(s.schedules)
}
func (s scheduleSorter) Less(i, j int) bool {
	return s.schedules[i].t.Before(s.schedules[j].t)
}
func (s scheduleSorter) Swap(i, j int) {
	schedule := s.schedules[i]
	s.schedules[i] = s.schedules[j]
	s.schedules[j] = schedule
}
func (c *Contest) sortSchedules() {
	sort.Stable(scheduleSorter{c})
}

func (c *Contest) startSchedule() {
	if len(c.schedules) == 0 {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.loaded {
		return
	}
	for len(c.schedules) > 0 && time.Now().After(c.schedules[0].t) {
		c.schedules[0].f()
		c.schedules = c.schedules[1:]
	}
	if len(c.schedules) == 0 {
		return
	}
	d := c.schedules[0].t.Sub(time.Now())
	c.scheduleTimer = time.AfterFunc(d, c.startSchedule)
}

// Call exactly once to unload the contest and save state into disk.
func (c *Contest) unload() {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.WithField("contestId", c.id).Info("Unloading contest")
	c.ranklist.Unload()
	for _, player := range c.players {
		c.unloadPlayer(player)
	}
	c.players = nil
	c.loaded = false
	c.cancelFunc()
	close(c.updateChan)
	close(c.playerUpdateChan)
	if c.scheduleTimer != nil {
		c.scheduleTimer.Stop()
	}
}

func handleWrites(ctx context.Context, ch chan mongo.WriteModel, coll *mongo.Collection, contestId primitive.ObjectID) {
	for {
		var writes []mongo.WriteModel
		select {
		case writeModel, ok := <-ch:
			if !ok {
				return
			}
			writes = append(writes, writeModel)
		loop:
			for {
				select {
				case writeModel, ok = <-ch:
					if !ok {
						break loop
					}
					writes = append(writes, writeModel)
				default:
					break loop
				}
			}
			if _, err := coll.BulkWrite(ctx, writes); err != nil {
				log.WithField("contestId", contestId).WithField("writeCount", len(writes)).Error("Failed to write to contest model: ", err)
				return
			} else {
				log.WithField("contestId", contestId).WithField("writeCount", len(writes)).Debug("Applied updates")
			}
		}
	}
}

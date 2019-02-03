package core

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type ContestPlayer struct {
	modelId  primitive.ObjectID
	userId   primitive.ObjectID
	problems map[string]*ContestPlayerProblem
}
type ContestPlayerProblem struct {
	subscriptions []*contestPlayerSubscription
}
type contestPlayerSubscription struct {
	c            *Contest
	userId       primitive.ObjectID
	submissionId primitive.ObjectID
	penaltyTime  time.Duration
	done         bool
	score        float64
}

func (s *contestPlayerSubscription) HandleNewScore(done bool, score float64) {
	s.c.lock.Lock()
	defer s.c.lock.Unlock()
	log.WithField("contestId", s.c.id).WithField("done", done).WithField("score", score).Info("Received submission score")
	s.done = done
	s.score = score

	p := s.c.players[s.userId]
	s.c.updatePlayerRankInfo(p)
}
func (c *Contest) updatePlayerRankInfo(player *ContestPlayer) {
	rankInfo := new(ContestPlayerRankInfo)
	rankInfo.problems = make(map[string]*ContestPlayerRankInfoProblem)
	for key, problem := range player.problems {
		problemInfo := new(ContestPlayerRankInfoProblem)
		for _, subscription := range problem.subscriptions {
			submissionInfo := &ContestPlayerRankInfoSubmission{
				Done:        subscription.done,
				Score:       subscription.score,
				PenaltyTime: subscription.penaltyTime,
			}
			problemInfo.submissions = append(problemInfo.submissions, submissionInfo)
		}
		rankInfo.problems[key] = problemInfo
	}
	c.ranklist.UpdatePlayer(player.userId, rankInfo)
}

func (c *Contest) loadPlayer(contestPlayerModel *model.ContestPlayer) {
	player := new(ContestPlayer)
	player.modelId = contestPlayerModel.Id
	player.userId = contestPlayerModel.User
	player.problems = make(map[string]*ContestPlayerProblem)
	for i, problemEntryModel := range contestPlayerModel.Problems {
		problemEntry := new(ContestPlayerProblem)
		for _, submission := range problemEntryModel.Submissions {
			subscription := &contestPlayerSubscription{
				c:            c,
				userId:       player.userId,
				submissionId: submission.SubmissionId,
				penaltyTime:  submission.PenaltyTime,
			}
			c.c.SubscribeSubmission(subscription.submissionId, subscription)
			problemEntry.subscriptions = append(problemEntry.subscriptions, subscription)
		}
		player.problems[i] = problemEntry
	}
	c.players[player.userId] = player
}

func (*Contest) unloadPlayer(p *ContestPlayer) {
	for _, problemEntry := range p.problems {
		for _, subscription := range problemEntry.subscriptions {
			subscription.c.c.UnsubscribeSubmission(subscription.submissionId, subscription)
		}
	}
}

package impl_leveldb

import (
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/syzoj/syzoj-ng-go/app/judge_traditional"
)

var log = logrus.StandardLogger()

type judgeService struct {
	db          *leveldb.DB
	judgeQueue  chan int64
	count       int64
	submissions sync.Map
	upgrader    websocket.Upgrader
	closeGroup  sync.WaitGroup
	closeChan   chan struct{}
	isClosed    int32
	clients     sync.Map
}

type submissionEntry struct {
	Tag       int64
	Language  string
	Code      string
	ProblemId uuid.UUID
	Callback  judge_traditional.Callback
}

func NewJudgeService(db *leveldb.DB) (judge_traditional.Service, error) {
	s := &judgeService{
		judgeQueue: make(chan int64, 1000),
		closeChan:  make(chan struct{}),
		upgrader:   websocket.Upgrader{},
		db:         db,
	}
	return s, nil
}

func (e *submissionEntry) getFields() logrus.Fields {
	return logrus.Fields{
		"Tag":       e.Tag,
		"ProblemId": e.ProblemId,
	}
}

func (ps *judgeService) QueueSubmission(sub *judge_traditional.Submission, callback judge_traditional.Callback) (judge_traditional.Task, error) {
	var id = atomic.AddInt64(&ps.count, 1)
	entry := &submissionEntry{
		Tag:       id,
		Language:  sub.Language,
		Code:      sub.Code,
		ProblemId: sub.ProblemId,
		Callback:  callback,
	}
	ps.submissions.Store(id, entry)
	select {
	case ps.judgeQueue <- id:
		log.WithFields(entry.getFields()).Debug("Queued submission")
		return nil, nil
	default:
		ps.submissions.Delete(id)
		return nil, judge_traditional.ErrQueueFull
	}
}

func (ps *judgeService) Close() error {
	atomic.StoreInt32(&ps.isClosed, 1)
	close(ps.closeChan)
	ps.closeGroup.Wait()
	// TODO: Save queue to LevelDB
	return nil
}

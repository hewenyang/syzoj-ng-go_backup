package service

import (
	"container/list"
	"context"
	"database/sql"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/judge/rpc"
	"github.com/syzoj/syzoj-ng-go/util"
)

func (s *serv) CreateSubmission(ctx context.Context, req *rpc.CreateSubmissionRequest) (*rpc.CreateSubmissionResponse, error) {
	id := util.RandomString(12)
	_, err := s.db.ExecContext(ctx, "INSERT INTO submissions (id, test_data, submit_data) VALUES (?, ?, ?)", id, util.ProtoSql{req.TestData}, util.ProtoSql{req.SubmitData})
	if err != nil {
		log.WithError(err).Error("Failed to update database")
		return nil, err
	}
	return &rpc.CreateSubmissionResponse{SubmissionId: proto.String(id)}, nil
}

func (s *serv) GetSubmission(ctx context.Context, req *rpc.GetSubmissionRequest) (*rpc.GetSubmissionResponse, error) {
	testData := &common.Data{}
	submitData := &common.Data{}
	result := &common.Data{}
	err := s.db.QueryRowContext(ctx, "SELECT test_data, submit_data, result FROM submissions WHERE id=?", req.GetSubmissionId()).Scan(util.ProtoSql{testData}, util.ProtoSql{submitData}, util.ProtoSql{result})
	if err != nil {
		if err == sql.ErrNoRows {
			return &rpc.GetSubmissionResponse{Error: rpc.Error_SubmissionNotFound.Enum()}, nil
		}
		log.WithError(err).Error("Failed to query database")
		return nil, err
	}
	return &rpc.GetSubmissionResponse{
		TestData:   testData,
		SubmitData: submitData,
		Result:     result,
	}, nil
}

func (s *serv) SubscribeSubmission(req *rpc.SubscribeSubmissionRequest, cli rpc.Judge_SubscribeSubmissionServer) error {
	ctx := cli.Context()
	submissionId := req.GetSubmissionId()
	var x int
	err := s.db.QueryRowContext(ctx, "SELECT 1 FROM submissions WHERE id=?", submissionId).Scan(&x)
	if err != nil {
		if err == sql.ErrNoRows {
			return cli.Send(&rpc.SubscribeSubmissionResponse{Error: rpc.Error_SubmissionNotFound.Enum()})
		}
		log.WithError(err).Error("Failed to query database")
		return err
	}

	l, w := s.subscribe(submissionId)
	defer s.unsubscribe(submissionId, l)
	for {
		select {
		case <-ctx.Done():
			break
		case d, ok := <-w:
			if !ok {
				break
			}
			err := cli.Send(&rpc.SubscribeSubmissionResponse{Result: d})
			if err != nil {
				break
			}
		}
	}
}

func (s *serv) HandleSubmission(cli rpc.Judge_HandleSubmissionServer) error {
	ctx := cli.Context()
	data, err := cli.Recv()
	if err != nil {
		return err
	}
	submissionId := data.GetSubmissionId()
	var x int
	err = s.db.QueryRowContext(ctx, "SELECT 1 FROM submissions WHERE id=?", submissionId).Scan(&x)
	if err != nil {
		if err == sql.ErrNoRows {
			return cli.SendAndClose(&rpc.HandleSubmissionResponse{Error: rpc.Error_SubmissionNotFound.Enum()})
		}
		log.WithError(err).Error("Failed to query database")
		return err
	}
	for {
		data, err := cli.Recv()
		if err != nil {
			break
		}
		s.publish(submissionId, data.Result, data.GetDone())
	}
	return nil
}

func (s *serv) subscribe(submissionId string) (*list.Element, <-chan *common.Data) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.subs[submissionId]
	if !ok {
		l = list.New()
		s.subs[submissionId] = l
	}
	ch := make(chan *common.Data)
	return l.PushBack(ch), ch
}

func (s *serv) unsubscribe(submissionId string, e *list.Element) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l := s.subs[submissionId]
	ch := l.Remove(e).(chan *common.Data)
	close(ch)
	if l.Len() == 0 {
		delete(s.subs, submissionId)
	}
}

func (s *serv) publish(submissionId string, d *common.Data, done bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.subs[submissionId]
	if !ok {
		return
	}
	for e := l.Front(); e != nil; e = e.Next() {
		ch := e.Value.(chan *common.Data)
		select {
		case ch <- d:
		default:
		}
		if done {
			close(ch)
		}
	}
	if done {
		l.Init()
		delete(s.subs, submissionId)
	}
}

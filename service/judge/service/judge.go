package service

import (
	"container/list"
	"context"
	"database/sql"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/judge/rpc"
	"github.com/syzoj/syzoj-ng-go/util"
)

type subEntry struct {
	l list.List
	v atomic.Value
}

func (s *serv) CreateSubmission(ctx context.Context, req *rpc.CreateSubmissionRequest) (*rpc.CreateSubmissionResponse, error) {
	id := util.RandomString(12)
	_, err := s.db.ExecContext(ctx, "INSERT INTO submissions (id, test_data, submit_data, done) VALUES (?, ?, ?, FALSE)", id, util.ProtoSql{req.TestData}, util.ProtoSql{req.SubmitData})
	if err != nil {
		s.log.WithError(err).Error("Failed to update database")
		return nil, err
	}
	return &rpc.CreateSubmissionResponse{SubmissionId: proto.String(id)}, nil
}

func (s *serv) GetSubmission(ctx context.Context, req *rpc.GetSubmissionRequest) (*rpc.GetSubmissionResponse, error) {
	testData := &common.Data{}
	submitData := &common.Data{}
	result := &common.Data{}
	var done bool
	err := s.db.QueryRowContext(ctx, "SELECT test_data, submit_data, result, done FROM submissions WHERE id=?", req.GetSubmissionId()).Scan(util.ProtoSql{testData}, util.ProtoSql{submitData}, util.ProtoSql{result}, &done)
	if err != nil {
		if err == sql.ErrNoRows {
			return &rpc.GetSubmissionResponse{Error: rpc.Error_SubmissionNotFound.Enum()}, nil
		}
		s.log.WithError(err).Error("Failed to query database")
		return nil, err
	}
	return &rpc.GetSubmissionResponse{
		TestData:   testData,
		SubmitData: submitData,
		Result:     result,
		Done:       proto.Bool(done),
	}, nil
}

func (s *serv) SubscribeSubmission(req *rpc.SubscribeSubmissionRequest, cli rpc.Judge_SubscribeSubmissionServer) error {
	ctx := cli.Context()
	submissionId := req.GetSubmissionId()
	var done bool
	err := s.db.QueryRowContext(ctx, "SELECT done FROM submissions WHERE id=?", submissionId).Scan(&done)
	if err != nil {
		if err == sql.ErrNoRows {
			return cli.Send(&rpc.SubscribeSubmissionResponse{Error: rpc.Error_SubmissionNotFound.Enum()})
		}
		s.log.WithError(err).Error("Failed to query database")
		return err
	}

	l, n, v := s.subscribe(submissionId)
	defer s.unsubscribe(submissionId, l)
	for {
		select {
		case <-ctx.Done():
			break
		case _, ok := <-n:
			if !ok {
				break
			}
			d := v.Load().(*common.Data)
			err := cli.Send(&rpc.SubscribeSubmissionResponse{Result: d})
			if err != nil {
				break
			}
		}
	}
}

// TODO: handle concurrency
func (s *serv) HandleSubmission(cli rpc.Judge_HandleSubmissionServer) error {
	ctx := cli.Context()
	data, err := cli.Recv()
	if err != nil {
		return err
	}
	submissionId := data.GetSubmissionId()
	var done bool
	result := &common.Data{}
	err = s.db.QueryRowContext(ctx, "SELECT done, result FROM submissions WHERE id=?", submissionId, util.ProtoSql{result}).Scan(&done)
	if err != nil {
		if err == sql.ErrNoRows {
			return cli.SendAndClose(&rpc.HandleSubmissionResponse{Error: rpc.Error_SubmissionNotFound.Enum()})
		}
		s.log.WithError(err).Error("Failed to query database")
		return err
	}
	if done {
		cli.SendAndClose(&rpc.HandleSubmissionResponse{Error: rpc.Error_SubmissionDone.Enum()})
		return nil
	}
	for {
		data, err := cli.Recv()
		if err != nil {
			break
		}
		s.publish(submissionId, data.Result, data.GetDone())
		if data.GetDone() {
			_, err := s.db.ExecContext(ctx, "UPDATE submissions SET result=?, done=TRUE WHERE id=?", util.ProtoSql{data.Result}, submissionId)
			if err != nil {
				s.log.WithError(err).Error("Failed to update database")
				return err
			}
			break
		}
	}
	cli.SendAndClose(&rpc.HandleSubmissionResponse{})
	return nil
}

func (s *serv) subscribe(submissionId string) (*list.Element, <-chan struct{}, *atomic.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.subs[submissionId]
	if !ok {
		e = &subEntry{}
		e.v.Store(&common.Data{})
		s.subs[submissionId] = e
	}
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return e.l.PushBack(ch), ch, &e.v
}

func (s *serv) unsubscribe(submissionId string, m *list.Element) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.subs[submissionId]
	e.l.Remove(m)
	if e.l.Len() == 0 {
		delete(s.subs, submissionId)
	}
}

func (s *serv) publish(submissionId string, d *common.Data, done bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.subs[submissionId]
	if !ok {
		return
	}
	e.v.Store(d)
	for m := e.l.Front(); m != nil; m = m.Next() {
		ch := m.Value.(chan struct{})
		select {
		case ch <- struct{}{}:
		default:
		}
		if done {
			close(ch)
		}
	}
	if done {
		delete(s.subs, submissionId)
	}
}

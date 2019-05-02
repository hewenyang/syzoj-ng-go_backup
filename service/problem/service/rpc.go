package service

import (
	"context"
	"database/sql"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/model/common"
	kafkapb "github.com/syzoj/syzoj-ng-go/service/problem/kafka"
	"github.com/syzoj/syzoj-ng-go/service/problem/rpc"
	"github.com/syzoj/syzoj-ng-go/util"
)

func (s *serv) CreateProblem(ctx context.Context, req *rpc.CreateProblemRequest) (*rpc.CreateProblemResponse, error) {
	title := req.GetTitle()
	statement := &common.ProblemStatement{}
	problemId := util.RandomString(12)
	_, err := s.db.ExecContext(ctx, "INSERT INTO problems (id, title, statement) VALUES (?, ?, ?)", problemId, title, util.ProtoSql{statement})
	if err != nil {
		s.log.WithError(err).Error("Failed to update database")
		return nil, err
	}
	return &rpc.CreateProblemResponse{ProblemId: proto.String(problemId)}, nil
}

func (s *serv) GetProblem(ctx context.Context, req *rpc.GetProblemRequest) (*rpc.GetProblemResponse, error) {
	problemId := req.GetProblemId()
	var title string
	statement := &common.ProblemStatement{}
	testData := &common.Data{}
	err := s.db.QueryRowContext(ctx, "SELECT title, statement, test_data FROM problems WHERE id=?", problemId).Scan(&title, util.ProtoSql{statement}, util.ProtoSql{testData})
	if err != nil {
		if err == sql.ErrNoRows {
			return &rpc.GetProblemResponse{Error: rpc.Error_ProblemNotFound.Enum()}, nil
		}
		s.log.WithError(err).Error("Failed to query database")
		return nil, err
	}
	return &rpc.GetProblemResponse{
		Title:     proto.String(title),
		Statement: statement,
		TestData:  testData,
	}, nil
}

func (s *serv) UpdateProblemStatement(ctx context.Context, req *rpc.UpdateProblemStatementRequest) (*rpc.UpdateProblemStatementResponse, error) {
	problemId, statement := req.GetProblemId(), req.GetStatement()
	if statement == nil {
		return &rpc.UpdateProblemStatementResponse{Error: rpc.Error_BadRequest.Enum()}, nil
	}
	res, err := s.db.ExecContext(ctx, "UPDATE problems SET statement=? WHERE id=?", util.ProtoSql{statement}, problemId)
	if err != nil {
		s.log.WithError(err).Error("Failed to update database")
		return nil, err
	} else if n, err := res.RowsAffected(); err == nil && n != 1 {
		return &rpc.UpdateProblemStatementResponse{Error: rpc.Error_ProblemNotFound.Enum()}, nil
	}
	return &rpc.UpdateProblemStatementResponse{}, nil
}

func (s *serv) UpdateProblemTestData(ctx context.Context, req *rpc.UpdateProblemTestDataRequest) (*rpc.UpdateProblemTestDataResponse, error) {
	problemId, testData := req.GetProblemId(), req.GetTestData()
	if testData == nil {
		return &rpc.UpdateProblemTestDataResponse{Error: rpc.Error_BadRequest.Enum()}, nil
	}
	res, err := s.db.ExecContext(ctx, "UPDATE problems SET test_data=? WHERE id=?", util.ProtoSql{testData}, problemId)
	if err != nil {
		s.log.WithError(err).Error("Failed to update database")
		return nil, err
	} else if n, err := res.RowsAffected(); err == nil && n != 1 {
		return &rpc.UpdateProblemTestDataResponse{Error: rpc.Error_ProblemNotFound.Enum()}, nil
	}
	return &rpc.UpdateProblemTestDataResponse{}, nil
}

func (s *serv) SubmitProblem(ctx context.Context, req *rpc.SubmitProblemRequest) (*rpc.SubmitProblemResponse, error) {
	problemId, submitData := req.GetProblemId(), req.GetSubmitData()
	if submitData == nil {
		return &rpc.SubmitProblemResponse{Error: rpc.Error_BadRequest.Enum()}, nil
	}
	testData := &common.Data{}
	err := s.db.QueryRowContext(ctx, "SELECT test_data FROM problems WHERE id=?", problemId).Scan(util.ProtoSql{testData})
	if err != nil {
		if err == sql.ErrNoRows {
			return &rpc.SubmitProblemResponse{Error: rpc.Error_ProblemNotFound.Enum()}, nil
		}
		s.log.WithError(err).Error("Failed to query database")
		return nil, err
	}
	submissionId, err := s.judgeCli.CreateSubmission(ctx, testData, submitData)
	if err != nil {
		return nil, err
	}
	s.writeKafkaMessage(&kafkapb.ProblemEvent{ProblemId: proto.String(problemId), Event: &kafkapb.ProblemEvent_Submission{Submission: &kafkapb.ProblemSubmissionEvent{SubmissionId: proto.String(submissionId)}}})
	return &rpc.SubmitProblemResponse{SubmissionId: proto.String(submissionId)}, nil
}

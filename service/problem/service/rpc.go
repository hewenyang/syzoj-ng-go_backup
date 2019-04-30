package service

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/problem/rpc"
	"github.com/syzoj/syzoj-ng-go/util"
)

func (s *serv) CreateProblem(ctx context.Context, req *rpc.CreateProblemRequest) (*rpc.CreateProblemResponse, error) {
	title := req.GetTitle()
	statement := &common.ProblemStatement{}
	problemId := util.RandomString(12)
	_, err := s.db.ExecContext(ctx, "INSERT INTO problems (id, title, statement) VALUES (?, ?, ?)", problemId, title, util.ProtoSql{statement})
	if err != nil {
		log.WithError(err).Error("Failed to update database")
		return nil, err
	}
	problem := &common.Problem{Id: proto.String(problemId)}
	return &rpc.CreateProblemResponse{Problem: problem}, nil
}

func (s *serv) ListProblem(ctx context.Context, req *rpc.ListProblemRequest) (*rpc.ListProblemResponse, error) {
	var problems []*common.Problem
	rows, err := s.db.QueryContext(ctx, "SELECT id, title, statement FROM problems")
	if err != nil {
		log.WithError(err).Error("Failed to query database")
		return nil, err
	}
	for rows.Next() {
		var problemId string
		var title string
		statement := &common.ProblemStatement{}
		if err := rows.Scan(&problemId, &title, util.ProtoSql{statement}); err != nil {
			log.WithError(err).Error("Failed to query database")
			return nil, err
		}
		problems = append(problems, &common.Problem{Id: proto.String(problemId), Title: proto.String(title), Statement: statement})
	}
	return &rpc.ListProblemResponse{Problem: problems}, nil
}

func (s *serv) UpdateProblemStatement(ctx context.Context, req *rpc.UpdateProblemStatementRequest) (*rpc.UpdateProblemStatementResponse, error) {
	problemId, statement := req.GetProblemId(), req.GetStatement()
	if statement == nil {
		return &rpc.UpdateProblemStatementResponse{Error: rpc.Error_BadRequest.Enum()}, nil
	}
	_, err := s.db.ExecContext(ctx, "UPDATE problems SET statement=? WHERE id=?", problemId, util.ProtoSql{statement})
	if err != nil {
		log.WithError(err).Error("Failed to update database")
		return nil, err
	}
	return &rpc.UpdateProblemStatementResponse{}, nil
}

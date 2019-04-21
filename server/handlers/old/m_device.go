package handlers

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
)

func populateProblemsetFull(ctx context.Context, ref database.ProblemsetRef) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	problemset, err := s.GetDB().GetProblemset(ctx, ref)
	if err != nil {
		log.WithError(err).Error("Failed to get problemset")
		return server.ErrBusy
	}
	res := &model.Problemset{}
	res.Id = proto.String(string(problemset.GetId()))
	res.Title = problemset.Title
	rows, err := s.GetDB().QueryContext(ctx, "SELECT id FROM problem WHERE problemset=?", ref)
	if err != nil {
		log.WithError(err).Error("Failed to get problemset entries")
		return server.ErrBusy
	}
	var problems []database.ProblemRef
	for rows.Next() {
		var pref database.ProblemRef
		if err := rows.Scan(&pref); err != nil {
			log.WithError(err).Error("Failed to get problem entry id")
			return server.ErrBusy
		}
		problems = append(problems, pref)
		res.Problem = append(res.Problem, string(pref))
	}
	if rows.Err() != nil {
		log.WithError(err).Error("Failed to get problemset entries")
		return server.ErrBusy
	}
	c.Mutate(&model.UpdateProblemset{Problemset: res})
	for _, problemRef := range problems {
		if err := populateProblemEntry(ctx, problemRef); err != nil {
			return err
		}
	}
	return nil
}

func populateProblemEntry(ctx context.Context, ref database.ProblemRef) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	problem, err := s.GetDB().GetProblem(ctx, ref)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	}
	if problem == nil {
		return server.ErrNotFound
	}
	res := &model.Problem{}
	res.Id = proto.String(string(problem.GetId()))
	// assume problem.Problem != nil
	res.Title = problem.Problem.Title
	c.Mutate(&model.UpdateProblem{Problem: res})
	return nil
}

func populateProblemFull(ctx context.Context, ref database.ProblemRef) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	problem, err := s.GetDB().GetProblem(ctx, ref)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	}
	if problem == nil {
		return server.ErrNotFound
	}
	res := &model.Problem{}
	res.Id = proto.String(string(problem.GetId()))
	// assume problem.Problem != nil
	res.Title = problem.Problem.Title
	res.Statement = problem.Problem.Statement
	res.Source = problem.Problem.Source
	res.Judge = problem.Problem.Judge
	c.Mutate(&model.UpdateProblem{Problem: res})
	return nil
}

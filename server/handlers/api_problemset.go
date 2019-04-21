package handlers

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
	"github.com/syzoj/syzoj-ng-go/server/device"
)

func Handle_Problemset_Create(ctx context.Context) error {
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		return server.ErrBusy
	}
	req := &model.ProblemsetCreateRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	if dev == nil || dev.User == nil {
		c.SendResult(&model.ProblemsetCreateResponse{
			Success: proto.Bool(false),
			Reason:  proto.String("Not logged in"),
		})
		return nil
	}
	problemset := &database.Problemset{}
	problemset.User = dev.User
	problemset.Problemset = &model.Problemset{}
	problemset.Problemset.Title = proto.String(req.GetProblemset().GetTitle())
	if err := s.GetDB().InsertProblemset(ctx, problemset); err != nil {
		log.WithError(err).Error("Failed to insert problemset")
		return server.ErrBusy
	}
	c.SendResult(&model.ProblemsetCreateResponse{
		Success: proto.Bool(true),
		Problemset: &model.Problemset{
			Id: proto.String(string(problemset.GetId())),
		},
	})
	return nil
}

func Get_Problemset(ctx context.Context) error {
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	vars := c.Vars()
	problemsetId := database.ProblemsetRef(vars["problemset_id"])
	pset, err := s.GetDB().GetProblemset(ctx, problemsetId)
	if err != nil {
		log.WithError(err).Error("Failed to get problemset")
		return server.ErrBusy
	}
	resp := &model.ProblemsetGetResponse{}
	resp.Problemset = &model.Problemset{}
	*resp.Problemset = *pset.Problemset
	resp.Problemset.Id = proto.String(string(pset.GetId()))
	rows, err := s.GetDB().QueryContext(ctx, "SELECT id FROM problem WHERE problemset = ?", problemsetId)
	if err != nil {
		log.WithError(err).Error("Failed to get problemset")
		return server.ErrBusy
	}
	for rows.Next() {
		var problemId database.ProblemRef
		if err := rows.Scan(&problemId); err != nil {
			log.WithError(err).Error("Failed to get problemset")
			return server.ErrBusy
		}
		problem, err := s.GetDB().GetProblem(ctx, problemId)
		if err != nil {
			log.WithError(err).Error("Failed to get problemset")
			return server.ErrBusy
		}
		problemEntry := &model.Problem{}
		problemEntry.Id = proto.String(string(problemId))
		problemEntry.Title = problem.Problem.Title
		resp.Problemset.Problem = append(resp.Problemset.Problem, problemEntry)
	}
	if rows.Err() != nil {
		log.WithError(err).Error("Failed to get problemset")
		return server.ErrBusy
	}
	c.SendResult(resp)
	return nil
}

func Handle_Problemset_Add_Problem(ctx context.Context) error {
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	vars := c.Vars()
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		return server.ErrBusy
	}
	req := &model.ProblemsetAddProblemRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	problemsetId := database.ProblemsetRef(vars["problemset_id"])
	pset, err := s.GetDB().GetProblemset(ctx, problemsetId)
	if err != nil {
		log.WithError(err).Error("Failed to get problemset")
		return server.ErrBusy
	}
	if dev == nil || dev.User == nil || dev.GetUser() != pset.GetUser() {
		c.SendResult(&model.ProblemsetAddProblemResponse{
			Success: proto.Bool(false),
			Reason:  proto.String("Permission denied"),
		})
		return nil
	}

	problem := &database.Problem{}
	problem.Problemset = database.CreateProblemsetRef(problemsetId)
	problem.Problem = &model.Problem{}
	problem.Problem.Title = req.Problem.Title
	problem.User = dev.User
	if err := s.GetDB().InsertProblem(ctx, problem); err != nil {
		log.WithError(err).Error("Failed to insert problem")
		return server.ErrBusy
	}
	c.SendResult(&model.ProblemsetAddProblemResponse{
		Success: proto.Bool(true),
		Problem: &model.Problem{
			Id: proto.String(string(problem.GetId())),
		},
	})
	return nil
}

func Get_Problem(ctx context.Context) error {
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	p, err := s.GetDB().GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problemset")
		return server.ErrBusy
	}
	resp := &model.ProblemGetResponse{}
	resp.Problem = &model.Problem{}
	*resp.Problem = *p.Problem
	resp.Problem.Id = proto.String(string(p.GetId()))
	c.SendResult(resp)
	return nil
}

func Handle_Problem_Update(ctx context.Context) error {
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		return server.ErrBusy
	}
	req := &model.ProblemUpdateRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	var ok bool = false
	_, err = s.GetDB().UpdateProblem(ctx, problemId, func(m *database.Problem) *database.Problem {
		if m.User == nil || dev.User == nil || m.GetUser() != dev.GetUser() {
			ok = false
			return m
		}
		m2 := &database.Problem{}
		*m2 = *m
		m2.Problem = &model.Problem{}
		*m2.Problem = *m.Problem
		p := req.GetProblem()
		if p != nil {
			if p.Title != nil {
				m2.Problem.Title = p.Title
			}
			if p.Statement != nil {
				m2.Problem.Statement = p.Statement
			}
			if p.Judge != nil {
				m2.Problem.Judge = p.Judge
			}
		}
		ok = true
		return m2
	})
	if err != nil {
		log.WithError(err).Error("Failed to update problem")
		return server.ErrBusy
	}
	c.SendResult(&model.ProblemsetAddProblemResponse{Success: proto.Bool(ok)})
	return nil
}

func init() {
	router.Action("/api/problemset/create", Handle_Problemset_Create)
	router.Get("/api/problemset/{problemset_id}", Get_Problemset)
	router.Action("/api/problemset/{problemset_id}/add-problem", Handle_Problemset_Add_Problem)
	router.Get("/api/problem/{problem_id}", Get_Problem)
	router.Action("/api/problem/{problem_id}/update", Handle_Problem_Update)
}

package handlers

import (
	"context"
	"io"

	"github.com/golang/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/model/judge"
	"github.com/syzoj/syzoj-ng-go/server"
	"github.com/syzoj/syzoj-ng-go/server/device"
)

func Handle_Problem_Remove_Judge(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemViewPage_RemoveJudgeRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	problem, err := s.GetDB().GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	}
	if problem == nil {
		return server.ErrNotFound
	}
	if problem.User == nil || dev.User == nil || problem.GetUser() != dev.GetUser() {
		return server.ErrPermissionDenied
	}

	var f error
	_, err = s.GetDB().UpdateProblem(ctx, problemId, func(m *database.Problem) *database.Problem {
		if m == nil {
			f = server.ErrNotFound
			return m
		} else {
			m2 := &database.Problem{}
			*m2 = *m
			m2.Problem = &model.Problem{}
			*m2.Problem = *m.Problem
			m2.Problem.Judge = nil
			return m2
		}
	})
	if err != nil {
		log.WithError(err).Error("Failed to update problem")
		return server.ErrBusy
	} else if f != nil {
		return f
	} else {
		return Get_Problem(ctx)
	}
}

func Handle_Problem_Edit_Statement(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemViewPage_EditStatementRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	problem, err := s.GetDB().GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	}
	if problem == nil {
		return server.ErrNotFound
	}
	if problem.User == nil || dev.User == nil || problem.GetUser() != dev.GetUser() {
		return server.ErrPermissionDenied
	}

	var f error
	_, err = s.GetDB().UpdateProblem(ctx, problemId, func(m *database.Problem) *database.Problem {
		if m == nil {
			f = server.ErrNotFound
			return m
		} else {
			m2 := &database.Problem{}
			*m2 = *m
			m2.Problem = &model.Problem{}
			*m2.Problem = *m.Problem
			m2.Problem.Statement = req.Statement
			return m2
		}
	})
	if err != nil {
		log.WithError(err).Error("Failed to update problem")
		return server.ErrBusy
	} else if f != nil {
		return f
	} else {
		return Get_Problem(ctx)
	}
}

func Handle_Problem_Add_Judge_Traditional(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemViewPage_AddJudgeTraditionalRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	problem, err := s.GetDB().GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	}
	if problem == nil {
		return server.ErrNotFound
	}
	if problem.User == nil || dev.User == nil || problem.GetUser() != dev.GetUser() {
		return server.ErrPermissionDenied
	}

	var f error
	_, err = s.GetDB().UpdateProblem(ctx, problemId, func(m *database.Problem) *database.Problem {
		if m == nil {
			f = server.ErrNotFound
			return m
		} else {
			m2 := &database.Problem{}
			*m2 = *m
			m2.Problem = &model.Problem{}
			*m2.Problem = *m.Problem
			if m2.Problem.Judge != nil {
				// Cannot update a problem with existing judge info
				f = server.ErrBusy
				return m
			}
			m2.Problem.Judge = &model.ProblemJudge{Judge: &model.ProblemJudge_Traditional{Traditional: req.Data}}
			return m2
		}
	})
	if err != nil {
		log.WithError(err).Error("Failed to update problem")
		return server.ErrBusy
	} else if f != nil {
		return f
	} else {
		return Get_Problem(ctx)
	}
}

func Handle_Problem_Submit_Judge_Traditional(ctx context.Context) error {
	c := server.GetWsApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemViewPage_SubmitJudgeTraditionalRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	dev, err := device.GetDeviceWs(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	problem, err := s.GetDB().GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	}
	if problem == nil {
		return server.ErrNotFound
	}
	jdata := problem.GetProblem().GetJudge()
	traditional, ok := jdata.GetJudge().(*model.ProblemJudge_Traditional)
	if !ok {
		log.WithError(err).Error("Problem judge is not traditional")
		return server.ErrBusy
	}
	jreq := &judge.JudgeRequest{
		Request: &judge.JudgeRequest_Traditional{Traditional: &judge.TraditionalJudgeRequest{
			ProblemId: proto.String(string(problemId)),
			Code:      req.Code,
			Data:      traditional.Traditional,
		}},
	}
	log.Info(jreq)
	j, err := s.GetJudge().JudgeSubmission(ctx, jreq)
	if err != nil {
		log.WithError(err).Error("Failed to judge submission")
		return server.ErrBusy
	}
	log.Info("Start")
	c.Mutate(&model.StartCurrentSubmission{})
	var result *judge.JudgeResponse
	for {
		result, err = j.GetResult(result)
		if err == io.EOF {
			break
		} else if err != nil {
			log.WithError(err).Error("Failed to get judge result")
			return server.ErrBusy
		}
		c.Mutate(&model.UpdateCurrentSubmission{Result: result})
	}
	c.Mutate(&model.StopCurrentSubmission{Result: result})
	return nil
}

func Get_Problem(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	if err := populateProblemFull(ctx, problemId); err != nil {
		return err
	}
	page := &model.ProblemViewPage{Problem: proto.String(string(problemId))}
	c.SendBody(page)
	return nil
}

func init() {
	router.Get("/api/problem/{problem_id:[0-9A-Za-z-_]{16}}", Get_Problem)
	router.Action("/api/problem/{problem_id:[0-9A-Za-z-_]{16}}/remove-judge", Handle_Problem_Remove_Judge)
	router.Action("/api/problem/{problem_id:[0-9A-Za-z-_]{16}}/edit-statement", Handle_Problem_Edit_Statement)
	router.Action("/api/problem/{problem_id:[0-9A-Za-z-_]{16}}/add-judge-traditional", Handle_Problem_Add_Judge_Traditional)
	router.WebSocket("/api/problem/{problem_id:[0-9A-Za-z-_]{16}}/judge/traditional/submit", Handle_Problem_Submit_Judge_Traditional)
}

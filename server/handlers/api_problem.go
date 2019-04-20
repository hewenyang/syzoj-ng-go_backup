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

func Handle_Problemset_Create_Problem(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	vars := c.Vars()
	problemsetRef := database.ProblemsetRef(vars["problemset_id"])
	req := &model.ProblemsetPage_CreateProblemRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	if req.ProblemTitle == nil {
		return server.ErrBadRequest
	}
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	problemset, err := s.GetDB().GetProblemset(ctx, problemsetRef)
	if err != nil {
		log.WithError(err).Error("Failed to get problemset")
		return server.ErrBusy
	}
	if problemset == nil {
		return server.ErrNotFound
	}
	if problemset.User == nil || dev.User == nil || problemset.GetUser() != dev.GetUser() {
		return server.ErrPermissionDenied
	}
	pb := &database.Problem{}
	pb.Title = req.ProblemTitle
	pb.User = dev.User
	pb.Problemset = problemset.Id
	pb.Problem = &model.Problem{}
	if err := s.GetDB().InsertProblem(ctx, pb); err != nil {
		log.WithError(err).Error("Failed to insert problem")
		return server.ErrBusy
	}
	c.Redirect("/problem/" + string(pb.GetId()))
	return nil
}

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
	j, err := s.GetJudge().JudgeSubmission(ctx, jreq)
	if err != nil {
		log.WithError(err).Error("Failed to judge submission")
		return server.ErrBusy
	}
	var result *judge.JudgeResponse
	for {
		result, err = j.GetResult(result)
		if err == io.EOF {
			break
		} else if err != nil {
			log.WithError(err).Error("Failed to get judge result")
			return server.ErrBusy
		}
		c.SendValue(result)
	}
	return nil
}

func Get_Problem(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	p, err := s.GetDB().GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	}
	if p == nil {
		return server.ErrNotFound
	}
	page := &model.ProblemViewPage{}
	page.ProblemTitle = p.Title
	page.ProblemStatement = p.Problem.Statement
	page.ProblemJudge = p.Problem.Judge
	c.SendBody(page)
	return nil
}

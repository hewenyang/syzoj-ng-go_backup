package handlers

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/model/judge"
	"github.com/syzoj/syzoj-ng-go/server"
	"github.com/syzoj/syzoj-ng-go/server/device"
)

func Get_Problems(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	page := new(model.ProblemsPage)
	rows, err := s.GetDB().QueryContext(ctx, "SELECT id FROM problem_entry")
	if err != nil {
		log.WithError(err).Error("Failed to get problem entries")
		return server.ErrBusy
	}
	for rows.Next() {
		var ref database.ProblemEntryRef
		if err := rows.Scan(&ref); err != nil {
			log.WithError(err).Error("Failed to get problem entry id")
			return server.ErrBusy
		}
		problemEntry, err := s.GetDB().GetProblemEntry(ctx, ref)
		if err != nil || problemEntry == nil {
			log.WithError(err).Error("Failed to get problem entry")
			return server.ErrBusy
		}
		problem, err := s.GetDB().GetProblem(ctx, problemEntry.GetProblem())
		if err != nil || problemEntry == nil {
			log.WithError(err).Error("Failed to get problem")
			return server.ErrBusy
		}
		page.ProblemEntry = append(page.ProblemEntry, &model.ProblemsPage_ProblemEntry{
			Id:           proto.String(string(ref)),
			ProblemId:    proto.String(string(problemEntry.GetProblem())),
			ProblemTitle: problem.Title,
		})
	}
	if rows.Err() != nil {
		log.WithError(err).Error("Failed to get problem entries")
		return server.ErrBusy
	}
	c.SendBody(page)
	return nil
}

func Get_Problem_Create(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		c.Redirect("/login")
		return nil
	}
	c.SendBody(&model.ProblemCreatePage{})
	return nil
}

func Handle_Problem_Create(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemCreatePage_CreateRequest{}
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
	pb := new(database.Problem)
	pb.Title = req.ProblemTitle
	pb.User = dev.User
	pb.Problem = &model.Problem{}
	if err := s.GetDB().InsertProblem(ctx, pb); err != nil {
		log.WithError(err).Error("Failed to insert problem")
		return server.ErrBusy
	}
	c.Redirect("/problem/" + string(pb.GetId()))
	return nil
}

func Handle_Problems_Add_Problem(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemsPage_AddProblemRequest{}
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
	problemId := database.ProblemRef(req.GetProblemId())
	problem, err := s.GetDB().GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	} else if problem == nil {
		return server.ErrNotFound
	}
	if problem.User == nil || dev.User == nil || problem.GetUser() != dev.GetUser() {
		return server.ErrPermissionDenied
	}
	entry := &database.ProblemEntry{}
	entry.Problem = database.CreateProblemRef(problemId)
	if err := s.GetDB().InsertProblemEntry(ctx, entry); err != nil {
		log.WithError(err).Error("Failed to insert problem entry")
		return server.ErrBusy
	}
	s.GetDB().FlushProblemEntry(ctx, entry.GetId())
	return Get_Problems(ctx)
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
	} else if problem == nil {
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
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemViewPage_SubmitJudgeTraditionalRequest{}
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
	} else if problem == nil {
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
	result, _ := j.GetResult()
	c.Mutate(fmt.Sprintf("/page/problem/%s/judge/traditional", string(problemId)), "setResult", result)
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
	page := &model.ProblemViewPage{}
	page.ProblemTitle = p.Title
	page.ProblemStatement = p.Problem.Statement
	page.ProblemSource = p.Problem.Source
	page.ProblemJudge = p.Problem.Judge
	{
		rows, err := s.GetDB().QueryContext(ctx, "SELECT id FROM problem_entry WHERE problem=?", problemId)
		if err != nil {
			log.WithError(err).Error("Failed to get problem entry")
			return server.ErrBusy
		}
		for rows.Next() {
			var ref database.ProblemEntryRef
			if err := rows.Scan(&ref); err != nil {
				log.WithError(err).Error("Failed to get problem entries")
				return server.ErrBusy
			}
			problemEntry, err := s.GetDB().GetProblemEntry(ctx, ref)
			if err != nil || problemEntry == nil {
				log.WithError(err).Error("Failed to get problem entry")
				return server.ErrBusy
			}
			page.ProblemEntry = append(page.ProblemEntry, &model.ProblemViewPage_ProblemEntry{Id: proto.String(string(ref))})
		}
		if rows.Err() != nil {
			log.WithError(err).Error("Failed to get problem entries")
			return server.ErrBusy
		}
	}
	c.SendBody(page)
	return nil
}

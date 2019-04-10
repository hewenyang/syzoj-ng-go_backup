package handlers

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
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
	if err := s.GetDB().InsertProblem(ctx, pb); err != nil {
		log.WithError(err).Error("Failed to insert problem")
		return server.ErrBusy
	}
	c.Redirect("/problem/" + string(pb.GetId()))
	return nil
}

func Handle_Problem_Set_Public(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemViewPage_SetPublicRequest{}
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
	entry := &database.ProblemEntry{}
	entry.Problem = database.CreateProblemRef(problemId)
	if err := s.GetDB().InsertProblemEntry(ctx, entry); err != nil {
		log.WithError(err).Error("Failed to insert problem entry")
		return server.ErrBusy
	}
	return Get_Problem(ctx)
}

func Handle_Problem_Add_Statement(ctx context.Context, req *model.ProblemViewPage_AddStatementRequest) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
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
	// TODO
	return Get_Problem(ctx)
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

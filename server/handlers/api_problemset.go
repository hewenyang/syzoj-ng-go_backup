package handlers

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
	"github.com/syzoj/syzoj-ng-go/server/device"
)

func Get_Problemset_Create(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		c.Redirect("/login")
		return nil
	}
	c.SendBody(&model.ProblemsetCreatePage{})
	return nil
}

func Handle_Problemset_Create(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	req := &model.ProblemsetCreatePage_CreateRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	if req.ProblemsetTitle == nil {
		return server.ErrBadRequest
	}
	dev, err := device.GetDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	p := &database.Problemset{}
	p.Title = req.ProblemsetTitle
	p.User = dev.User
	if err := s.GetDB().InsertProblemset(ctx, p); err != nil {
		log.WithError(err).Error("Failed to insert problemset")
		return server.ErrBusy
	}
	c.Redirect("/problemset/" + string(p.GetId()))
	return nil
}

func Get_Problemset(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	vars := c.Vars()
	problemsetId := database.ProblemsetRef(vars["problemset_id"])
	page := &model.ProblemsetPage{}
	problemset, err := s.GetDB().GetProblemset(ctx, problemsetId)
	if err != nil {
		log.WithError(err).Error("Failed to get problemset")
		return server.ErrBusy
	}
	if problemset == nil {
		return server.ErrNotFound
	}
	rows, err := s.GetDB().QueryContext(ctx, "SELECT id FROM problem WHERE problemset=?", problemsetId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem entries")
		return server.ErrBusy
	}
	for rows.Next() {
		var ref database.ProblemRef
		if err := rows.Scan(&ref); err != nil {
			log.WithError(err).Error("Failed to get problem entry id")
			return server.ErrBusy
		}
		problem, err := s.GetDB().GetProblem(ctx, ref)
		if problem == nil || err != nil {
			log.WithError(err).Error("Failed to get problem")
			return server.ErrBusy
		}
		page.ProblemEntry = append(page.ProblemEntry, &model.ProblemsetPage_ProblemEntry{
			Id:           proto.String(string(ref)),
			ProblemId:    proto.String(string(problem.GetId())),
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

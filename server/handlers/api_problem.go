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
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return server.ErrBusy
	}
	defer txn.Rollback()
    page := new(model.ProblemsPage)
	var problemEntryRef []database.ProblemEntryRef
	rows, err := txn.QueryContext(ctx, "SELECT id FROM problem_entry")
	if err != nil {
		log.WithError(err).Error("Failed to get problem entries")
		return server.ErrBusy
	}
	if err := database.ScanAll(rows, &problemEntryRef); err != nil {
		log.WithError(err).Error("Failed to get problem entries")
		return server.ErrBusy
	}
	for _, ref := range problemEntryRef {
		problemEntry, err := txn.GetProblemEntry(ctx, ref)
		if err != nil || problemEntry == nil {
			log.WithError(err).Error("Failed to get problem entry")
			return server.ErrBusy
		}
		problem, err := txn.GetProblem(ctx, problemEntry.GetProblem())
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
	c.SendBody(page)
	return nil
}

func Get_Problem_Create(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return server.ErrBusy
	}
	defer txn.Rollback()
	dev, err := device.GetDevice(ctx, txn)
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

func Handle_Problem_Create(ctx context.Context, req *model.ProblemCreatePage_CreateRequest) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	if req.ProblemTitle == nil {
		return server.ErrBadRequest
	}
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return server.ErrBusy
	}
	defer txn.Rollback()
	dev, err := device.GetDevice(ctx, txn)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	pb := new(database.Problem)
	pb.Title = req.ProblemTitle
	pb.User = dev.User
	if err := txn.InsertProblem(ctx, pb); err != nil {
		log.WithError(err).Error("Failed to insert problem")
		return server.ErrBusy
	}
	if err := txn.Commit(ctx); err != nil {
		log.WithError(err).Error("Failed to commit transaction")
		return server.ErrBusy
	}
	c.Redirect("/problem/" + string(pb.GetId()))
	return nil
}

func Handle_Problem_Set_Public(ctx context.Context, req *model.ProblemViewPage_SetPublicRequest) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return server.ErrBusy
	}
	defer txn.Rollback()
	dev, err := device.GetDevice(ctx, txn)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	problem, err := txn.GetProblem(ctx, problemId)
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
	if err := txn.InsertProblemEntry(ctx, entry); err != nil {
		log.WithError(err).Error("Failed to insert problem entry")
		return server.ErrBusy
	}
	if err := txn.Commit(ctx); err != nil {
		log.WithError(err).Error("Failed to commit transaction")
		return server.ErrBusy
	}
	return Get_Problem(ctx)
}

func Handle_Problem_Add_Statement(ctx context.Context, req *model.ProblemViewPage_AddStatementRequest) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return server.ErrBusy
	}
	defer txn.Rollback()
	dev, err := device.GetDevice(ctx, txn)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to find device")
		return server.ErrBusy
	} else if err == device.ErrDeviceNotFound || dev.User == nil {
		return server.ErrNotLoggedIn
	}
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	problem, err := txn.GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	} else if problem == nil {
		return server.ErrNotFound
	}
	if problem.User == nil || dev.User == nil || problem.GetUser() != dev.GetUser() {
		return server.ErrPermissionDenied
	}
	statement := new(database.ProblemStatement)
	statement.User = database.CreateUserRef(dev.GetUser())
	statement.Problem = database.CreateProblemRef(problemId)
	statement.Data = req.Statement
	if err := txn.InsertProblemStatement(ctx, statement); err != nil {
		log.WithError(err).Error("Failed to insert problem statement")
		return server.ErrBusy
	}
	if err := txn.Commit(ctx); err != nil {
		log.WithError(err).Error("Failed to commit transaction")
		return server.ErrBusy
	}
	return Get_Problem(ctx)
}

func Get_Problem(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	s := server.GetServer(ctx)
	vars := c.Vars()
	problemId := database.ProblemRef(vars["problem_id"])
	txn, err := s.GetDB().OpenReadonlyTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return server.ErrBusy
	}
	defer txn.Rollback()
	p, err := txn.GetProblem(ctx, problemId)
	if err != nil {
		log.WithError(err).Error("Failed to get problem")
		return server.ErrBusy
	}
	page := &model.ProblemViewPage{}
	page.ProblemTitle = p.Title
	{
		var problemStatementRef []database.ProblemStatementRef
		rows, err := txn.QueryContext(ctx, "SELECT id FROM problem_statement WHERE problem=?", problemId)
		if err != nil {
			log.WithError(err).Error("Failed to get problem statements")
			return server.ErrBusy
		}
		if err := database.ScanAll(rows, &problemStatementRef); err != nil {
			log.WithError(err).Error("Failed to get problem statements")
			return server.ErrBusy
		}
		for _, ref := range problemStatementRef {
			problemStatement, err := txn.GetProblemStatement(ctx, ref)
			if err != nil || problemStatement == nil {
				log.WithError(err).Error("Failed to get problem statement")
				return server.ErrBusy
			}
			page.ProblemStatement = append(page.ProblemStatement, &model.ProblemViewPage_ProblemStatementEntry{Id: proto.String(string(ref)), Statement: problemStatement.Data})
		}
	}
	{
		var problemEntryRef []database.ProblemEntryRef
		rows, err := txn.QueryContext(ctx, "SELECT id FROM problem_entry WHERE problem=?", problemId)
		if err != nil {
			log.WithError(err).Error("Failed to get problem entry")
			return server.ErrBusy
		}
		if err := database.ScanAll(rows, &problemEntryRef); err != nil {
			log.WithError(err).Error("Failed to get problem entry")
			return server.ErrBusy
		}
		for _, ref := range problemEntryRef {
			problemEntry, err := txn.GetProblemEntry(ctx, ref)
			if err != nil || problemEntry == nil {
				log.WithError(err).Error("Failed to get problem entry")
				return server.ErrBusy
			}
			page.ProblemEntry = append(page.ProblemEntry, &model.ProblemViewPage_ProblemEntry{Id: proto.String(string(ref))})
		}
	}
	c.SendBody(page)
	return nil
}

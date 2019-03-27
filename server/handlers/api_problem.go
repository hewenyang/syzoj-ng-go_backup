package handlers

import (
    "context"

    "github.com/syzoj/syzoj-ng-go/server"
    "github.com/syzoj/syzoj-ng-go/model"
    "github.com/syzoj/syzoj-ng-go/database"
)

func Get_Problem_Create(ctx context.Context) error {
    c := server.GetApiContext(ctx)
    c.SendBody(&model.ProblemCreatePage{})
    return nil
}

func Handle_Problem_Create(ctx context.Context, req *model.ProblemCreateRequest) error {
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
    pb := new(database.Problem)
    pb.Title = req.ProblemTitle
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
    c.SendBody(&model.ProblemViewPage{
        ProblemTitle: p.Title,
    })
    return nil
}

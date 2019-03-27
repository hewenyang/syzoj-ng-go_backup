package handlers

import (
    "context"

    "github.com/syzoj/syzoj-ng-go/server"
    "github.com/syzoj/syzoj-ng-go/model"
)

func Get_Problem_Create(ctx context.Context) error {
    c := server.GetApiContext(ctx)
    c.SendBody(&model.ProblemCreatePage{})
    return nil
}

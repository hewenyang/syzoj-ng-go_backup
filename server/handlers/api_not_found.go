package handlers

import (
	"context"

	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
)

func Handle_Not_Found(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	c.SendBody(&model.NotFoundPage{})
	return nil
}

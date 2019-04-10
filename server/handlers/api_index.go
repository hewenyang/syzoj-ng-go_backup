package handlers

import (
	"context"

	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
)

func Get_Index(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	c.SendBody(&model.IndexPage{})
	return nil
}

package service

import (
	"context"

	"github.com/syzoj/syzoj-ng-go/service/api/model"
)

func (s *apiService) handleProblems(ctx context.Context, c *ApiContext) error {
	resp := &model.ProblemsResponse{}
	c.SendBody(resp)
	return nil
}

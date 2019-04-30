package service

import (
	"context"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/api/model"
)

func convProblem(p *common.Problem) *model.Problem {
	return &model.Problem{
		Id:        p.Id,
		Title:     p.Title,
		Statement: p.Statement,
	}
}

func (s *apiService) handleProblems(ctx context.Context, c *ApiContext) error {
	res, err := s.problemClient.ListProblem(ctx)
	if err != nil {
		return err
	}
	resp := &model.ProblemsResponse{}
	for _, problem := range res {
		resp.Problems = append(resp.Problems, convProblem(problem))
	}
	c.SendBody(resp)
	return nil
}

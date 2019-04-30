package client

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/problem/rpc"
)

type ProblemError rpc.Error

func (e ProblemError) Error() string {
	return rpc.Error(e).String()
}

func (c *Client) CreateProblem(ctx context.Context, title string) (*common.Problem, error) {
	resp, err := c.c.CreateProblem(ctx, &rpc.CreateProblemRequest{Title: proto.String(title)})
	if err != nil {
		return nil, err
	} else if resp.Error != nil {
		return nil, ProblemError(resp.GetError())
	}
	return resp.Problem, nil
}
func (c *Client) UpdateProblemStatement(ctx context.Context, problemId string, statement *common.ProblemStatement) error {
	resp, err := c.c.UpdateProblemStatement(ctx, &rpc.UpdateProblemStatementRequest{ProblemId: proto.String(problemId), Statement: statement})
	if err != nil {
		return err
	} else if resp.Error != nil {
		return ProblemError(resp.GetError())
	}
	return nil
}

func (c *Client) ListProblem(ctx context.Context) ([]*common.Problem, error) {
	resp, err := c.c.ListProblem(ctx, &rpc.ListProblemRequest{})
	if err != nil {
		return nil, err
	} else if resp.Error != nil {
		return nil, ProblemError(resp.GetError())
	}
	return resp.Problem, nil
}

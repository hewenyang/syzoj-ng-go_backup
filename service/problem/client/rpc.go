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

func (c *Client) CreateProblem(ctx context.Context, title string) (string, error) {
	resp, err := c.c.CreateProblem(ctx, &rpc.CreateProblemRequest{Title: proto.String(title)})
	if err != nil {
		return "", err
	} else if resp.Error != nil {
		return "", ProblemError(resp.GetError())
	}
	return resp.GetProblemId(), nil
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

package client

import (
	"context"
	"io"

	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/judge/rpc"
)

var log = logrus.StandardLogger()

type JudgeError rpc.Error

func (e JudgeError) Error() string {
	return rpc.Error(e).String()
}

func (c *Client) CreateSubmission(ctx context.Context, testData *common.Data, submitData *common.Data) (string, error) {
	res, err := c.c.CreateSubmission(ctx, &rpc.CreateSubmissionRequest{
		TestData:   testData,
		SubmitData: submitData,
	})
	if err != nil {
		return "", err
	} else if res.Error != nil {
		return "", JudgeError(res.GetError())
	} else {
		return res.GetSubmissionId(), nil
	}
}

func (c *Client) GetSubmission(ctx context.Context, submissionId string) (*common.Data, *common.Data, *common.Data, error) {
	res, err := c.c.GetSubmission(ctx, &rpc.GetSubmissionRequest{SubmissionId: proto.String(submissionId)})
	if err != nil {
		return nil, nil, nil, err
	} else if res.Error != nil {
		return nil, nil, nil, JudgeError(res.GetError())
	} else {
		return res.TestData, res.SubmitData, res.Result, nil
	}
}

type subscribeType = struct {
	Data *common.Data
	Err  error
}

// Use context for cancellation
func (c *Client) SubscribeSubmission(ctx context.Context, submissionId string) (<-chan subscribeType, error) {
	x, err := c.c.SubscribeSubmission(ctx, &rpc.SubscribeSubmissionRequest{SubmissionId: proto.String(submissionId)})
	if err != nil {
		return nil, err
	}
	ch := make(chan subscribeType)
	go func() {
		for {
			resp, err := x.Recv()
			if err != nil {
				if err != io.EOF {
					log.WithError(err).Error("Failed to receive from gRPC stream")
				}
				break
			}
			if resp.Error != nil {
				err = JudgeError(resp.GetError())
			}
			ch <- subscribeType{Data: resp.Result, Err: err}
		}
		close(ch)
	}()
	return ch, nil
}

func (c *Client) HandleSubmissionRequest(ctx context.Context, submissionId string, data <-chan struct {
	Data *common.Data
	Done bool
}) error {
	cli, err := c.c.HandleSubmission(ctx)
	if err != nil {
		return err
	}
	if err = cli.Send(&rpc.HandleSubmissionRequest{SubmissionId: proto.String(submissionId)}); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			break
		case d, ok := <-data:
			if !ok {
				break
			}
			err := cli.Send(&rpc.HandleSubmissionRequest{Result: d.Data, Done: proto.Bool(d.Done)})
			if err != nil {
				break
			}
		}
	}
	resp, err := cli.CloseAndRecv()
	if err != nil {
		return err
	} else if resp.Error != nil {
		return JudgeError(resp.GetError())
	}
	return nil
}

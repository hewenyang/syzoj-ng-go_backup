package client

import (
	"google.golang.org/grpc"

	"github.com/syzoj/syzoj-ng-go/fakenet"
	"github.com/syzoj/syzoj-ng-go/service/problem/rpc"
)

type Client struct {
	g *grpc.ClientConn
	c rpc.ProblemClient
}

func NewProblemClient() (*Client, error) {
	g, err := grpc.Dial("service-problem", grpc.WithInsecure(), grpc.WithContextDialer(fakenet.Base.DialContext))
	if err != nil {
		return nil, err
	}
	return &Client{g: g, c: rpc.NewProblemClient(g)}, nil
}

func (c *Client) Close() error {
	return c.g.Close()
}

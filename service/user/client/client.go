package client

import (
	"google.golang.org/grpc"

	"github.com/syzoj/syzoj-ng-go/fakenet"
	"github.com/syzoj/syzoj-ng-go/service/user/rpc"
)

type Client struct {
	g *grpc.ClientConn
	c rpc.UserClient
}

func NewUserClient() (*Client, error) {
	g, err := grpc.Dial("service-user", grpc.WithInsecure(), grpc.WithContextDialer(fakenet.Base.DialContext))
	if err != nil {
		return nil, err
	}
	return &Client{g: g, c: rpc.NewUserClient(g)}, nil
}

func (c *Client) Close() error {
	return c.g.Close()
}

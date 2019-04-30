package client

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/user/rpc"
)

type Error rpc.Error

var ErrDuplicateUserName = Error(rpc.Error_DuplicateUserName)
var ErrUserNotFound = Error(rpc.Error_UserNotFound)
var ErrPasswordIncorrect = Error(rpc.Error_PasswordIncorrect)

func (e Error) Error() string {
	return rpc.Error(e).String()
}

func (c *Client) RegisterUser(ctx context.Context, userName string, password string) (string, error) {
	res, err := c.c.RegisterUser(ctx, &rpc.RegisterUserRequest{UserName: proto.String(userName), Password: proto.String(password)})
	if err != nil {
		return "", err
	} else if res.Error != nil {
		return "", Error(res.GetError())
	}
	return res.GetUserId(), nil
}

func (c *Client) LoginUser(ctx context.Context, userName string, password string, info *common.DeviceInfo) (string, error) {
	res, err := c.c.LoginUser(ctx, &rpc.LoginUserRequest{UserName: proto.String(userName), Password: proto.String(password), DeviceInfo: info})
	if err != nil {
		return "", err
	} else if res.Error != nil {
		return "", Error(res.GetError())
	}
	return res.GetToken(), nil
}

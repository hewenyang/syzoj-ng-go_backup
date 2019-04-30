package service

import (
	"context"
	"net/http"
	"time"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/api/model"
	"github.com/syzoj/syzoj-ng-go/service/user/client"
)

func (s *apiService) handleLogin(ctx context.Context, c *ApiContext) error {
	body := &model.LoginRequest{}
	if err := c.ReadBody(body); err != nil {
		return ErrBadRequest
	}
	if body.UserName == nil || body.Password == nil {
		return ErrBadRequest
	}
	deviceId, err := s.userClient.LoginUser(ctx, *body.UserName, *body.Password, &common.DeviceInfo{})
	switch err {
	case client.ErrUserNotFound:
		return ErrUserNotFound
	case nil:
	default:
		return err
	}
	c.SetCookie(&http.Cookie{
		Name:    "SYZOJSESSION",
		Value:   deviceId,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24 * 365),
	})
	c.SendBody(&model.LoginResponse{})
	return nil
}

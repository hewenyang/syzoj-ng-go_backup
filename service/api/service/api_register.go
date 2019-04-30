package service

import (
	"context"
	"net/http"
	"time"

	"github.com/syzoj/syzoj-ng-go/service/api/model"
	"github.com/syzoj/syzoj-ng-go/service/user/client"
)

func (s *apiService) handleRegister(ctx context.Context, c *ApiContext) error {
	body := &model.RegisterRequest{}
	if err := c.ReadBody(body); err != nil {
		return ErrBadRequest
	}
	userName, password := body.GetUserName(), body.GetPassword()
	deviceId, err := s.userClient.RegisterUser(ctx, userName, password)
	switch err {
	case client.ErrDuplicateUserName:
		return ErrDuplicateUserName
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
	c.SendBody(&model.RegisterResponse{})
	return nil
}

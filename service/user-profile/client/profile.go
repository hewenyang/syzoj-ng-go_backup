package client

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/user-profile/rpc"
)

type UserProfileError rpc.Error

var ErrBadRequest = UserProfileError(rpc.Error_BadRequest)

func (e UserProfileError) Error() string {
	return rpc.Error(e).String()
}

func (c *Client) GetUserProfile(ctx context.Context, userId string) (*common.UserProfile, error) {
	if userId == "" {
		return nil, ErrBadRequest
	}
	res, err := c.c.GetProfile(ctx, &rpc.GetProfileRequest{UserId: proto.String(userId)})
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, UserProfileError(res.GetError())
	}
	return res.Profile, nil
}

func (c *Client) UpdateUserProfile(ctx context.Context, userId string, profile *common.UserProfile) error {
	if userId == "" || profile == nil {
		return ErrBadRequest
	}
	res, err := c.c.UpdateProfile(ctx, &rpc.UpdateProfileRequest{UserId: proto.String(userId), Profile: profile})
	if err != nil {
		return err
	}
	if res.Error != nil {
		return UserProfileError(res.GetError())
	}
	return nil
}

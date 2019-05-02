package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/syzoj/syzoj-ng-go/model/common"
	"github.com/syzoj/syzoj-ng-go/service/user-profile/rpc"
	"github.com/syzoj/syzoj-ng-go/util"
)

var ErrBadRequest = errors.New("Bad request")

func (s *serv) GetProfile(ctx context.Context, req *rpc.GetProfileRequest) (*rpc.GetProfileResponse, error) {
	profile := &common.UserProfile{}
	err := s.db.QueryRowContext(ctx, "SELECT profile FROM user_profile WHERE id=?", req.GetUserId()).Scan(util.ProtoSql{profile})
	if err != nil {
		s.log.WithError(err).Error("Failed to query database")
		return nil, err
	}
	if err == sql.ErrNoRows {
		return &rpc.GetProfileResponse{Error: rpc.Error_UserNotFound.Enum()}, nil
	}
	return &rpc.GetProfileResponse{Profile: profile}, nil
}

func (s *serv) UpdateProfile(ctx context.Context, req *rpc.UpdateProfileRequest) (*rpc.UpdateProfileResponse, error) {
	if req.Profile == nil {
		return nil, ErrBadRequest
	}
	_, err := s.db.ExecContext(ctx, "UPDATE user_profile SET profile=? WHERE id=?", util.ProtoSql{req.Profile}, req.GetUserId())
	if err != nil {
		s.log.WithError(err).Error("Failed to update database")
		return nil, err
	}
	return &rpc.UpdateProfileResponse{}, nil
}

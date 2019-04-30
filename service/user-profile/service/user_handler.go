package service

import (
	"context"

	"github.com/syzoj/syzoj-ng-go/model/common"
	userkpb "github.com/syzoj/syzoj-ng-go/service/user/kafka"
	"github.com/syzoj/syzoj-ng-go/util"
)

func (s *serv) handleUserEvent(ctx context.Context, msg *userkpb.UserEvent) error {
	switch e := msg.Event.(type) {
	case *userkpb.UserEvent_Register:
		_ = e
		return s.createUser(ctx, msg.GetUserId())
	}
	return nil
}

func (s *serv) createUser(ctx context.Context, userId string) error {
	profile := &common.UserProfile{}
	_, err := s.db.ExecContext(ctx, "INSERT INTO user_profile (id, profile) VALUES (?, ?)", userId, util.ProtoSql{profile})
	return err
}

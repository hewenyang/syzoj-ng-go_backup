package handlers

import (
	"context"
	"database/sql"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
	"github.com/syzoj/syzoj-ng-go/server/device"
)

func Handle_Login(ctx context.Context) error {
	var err error
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	req := &model.LoginRequest{}
	if err := c.ReadBody(req); err != nil {
		return server.ErrBadRequest
	}
	if dev, err := device.GetDevice(ctx); err != device.ErrDeviceNotFound {
		if err != nil {
			log.WithError(err).Error("Failed to find device")
			return server.ErrBusy
		}
		if dev.User != nil {
			c.SendResult(&model.LoginResponse{
				Success: proto.Bool(false),
				Reason:  proto.String("Already logged in"),
			})
			return nil
		}
	}
	var userRef database.UserRef
	if err = s.GetDB().QueryRowContext(ctx, "SELECT id FROM user WHERE user_name=?", req.GetUserName()).Scan(&userRef); err != nil {
		if err == sql.ErrNoRows {
			c.SendResult(&model.LoginResponse{
				Success: proto.Bool(false),
				Reason:  proto.String("User not found"),
			})
			return nil
		}
		log.WithError(err).Error("Handle_Login query failed")
		return server.ErrBusy
	}
	var user *database.User
	if user, err = s.GetDB().GetUser(ctx, userRef); err != nil || user == nil {
		log.WithError(err).Error("Handle_Login query failed")
		return server.ErrBusy
	}
	if user.Auth == nil {
		log.Warning("Handle_Login: user.Auth is nil")
		return server.ErrBusy
	}
	if bcrypt.CompareHashAndPassword(user.Auth.PasswordHash, []byte(req.GetPassword())) != nil {
		c.SendResult(&model.LoginResponse{
			Success: proto.Bool(false),
			Reason:  proto.String("Password incorrect"),
		})
		return nil
	}
	dev, err := device.NewDevice(ctx)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to create device")
		return server.ErrBusy
	} else {
		err = nil
	}
	if _, err = s.GetDB().UpdateDevice(ctx, dev.GetId(), func(dev *database.Device) *database.Device {
		dev2 := &database.Device{}
		*dev2 = *dev
		dev2.User = database.CreateUserRef(userRef)
		return dev2
	}); err != nil {
		log.WithError(err).Error("Failed to update device")
		return server.ErrBusy
	}
	c.SendResult(&model.LoginResponse{
		Success: proto.Bool(true),
	})
	return nil
}

func init() {
	router.Action("/api/login", Handle_Login)
}

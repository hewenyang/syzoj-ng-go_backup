package handlers

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
	"github.com/syzoj/syzoj-ng-go/server/device"
)

func Get_Login(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	c.SendBody(&model.LoginPage{})
	return nil
}

func Handle_Login(ctx context.Context, req *model.LoginRequest) error {
	var err error
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return server.ErrBusy
	}
	defer txn.Rollback()
	if dev, err := device.GetDevice(ctx, txn); err != device.ErrDeviceNotFound {
		if err != nil {
			log.WithError(err).Error("Failed to find device")
			return server.ErrBusy
		}
		if dev.User != nil {
			return server.ErrAlreadyLoggedIn
		}
	}
	var userRef database.UserRef
	if err = txn.QueryRowContext(ctx, "SELECT id FROM user WHERE user_name=?", req.GetUserName()).Scan(&userRef); err != nil {
		if err == sql.ErrNoRows {
			return server.ErrUserNotFound
		}
		log.WithError(err).Error("Handle_Login query failed")
		return server.ErrBusy
	}
	var user *database.User
	if user, err = txn.GetUser(ctx, userRef); err != nil || user == nil {
		log.WithError(err).Error("Handle_Login query failed")
		return server.ErrBusy
	}
	if user.Auth == nil {
		log.Warning("Handle_Login: user.Auth is nil")
		return server.ErrBusy
	}
	if bcrypt.CompareHashAndPassword(user.Auth.PasswordHash, []byte(req.GetPassword())) != nil {
		return server.ErrPasswordIncorrect
	}
	dev, err := device.NewDevice(ctx, txn)
	if err != nil && err != device.ErrDeviceNotFound {
		log.WithError(err).Error("Failed to create device")
		return server.ErrBusy
	} else {
		err = nil
	}
	dev.User = database.CreateUserRef(userRef)
	if err = txn.UpdateDevice(ctx, dev.GetId(), dev); err != nil {
		log.WithError(err).Error("Failed to update device")
		return server.ErrBusy
	}
	if err = txn.Commit(ctx); err != nil {
		log.WithError(err).Error("Failed to commit transaction")
		return server.ErrBusy
	}
	c.Redirect("/")
	return nil
}

package handlers

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/server"
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
	c.Redirect("/")
	return nil
}

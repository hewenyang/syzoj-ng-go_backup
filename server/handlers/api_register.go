package handlers

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
)

func Get_Register(ctx context.Context) error {
	c := server.GetApiContext(ctx)
	c.SendBody(&model.RegisterPage{})
	return nil
}

func Handle_Register(ctx context.Context, req *model.RegisterPage_RegisterRequest) error {
	var err error
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return server.ErrBusy
	}
	defer txn.Rollback()
	user := new(database.User)
	if req.UserName == nil || !model.CheckUserName(req.GetUserName()) {
		return server.ErrBadRequest
	}
	user.UserName = req.UserName
	user.Auth = new(model.UserAuth)
	if req.Password == nil {
		return server.ErrBadRequest
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 0)
	if err != nil {
		log.WithError(err).Error("Failed to generate passowrd")
		return server.ErrBusy
	}
	user.Auth.PasswordHash = passwordHash
	if err = txn.InsertUser(ctx, user); err != nil {
		log.WithError(err).Error("Handle_Register query failed")
		return server.ErrBusy
	}
	if err = txn.Commit(ctx); err != nil {
		log.WithError(err).Error("Handle_Register query failed")
		return server.ErrBusy
	}
	c.Redirect("/login")
	return nil
}

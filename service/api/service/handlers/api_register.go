package handlers

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/gogo/protobuf/proto"
	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
)

func Handle_Register(ctx context.Context) error {
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	req := &model.RegisterRequest{}
	if err := c.ReadBody(req); err != nil {
		return err
	}
	if req = req.Validate(); req == nil {
		return server.ErrBadRequest
	}
	if req.UserName == nil || !model.CheckUserName(req.GetUserName()) || req.Password == nil {
		return server.ErrBadRequest
	}

	// Lock username
	userName := req.GetUserName()
	o := s.GetOracle()
	o.Lock()
	for {
		_, found := o.Map[server.UserNameKey(userName)]
		if !found {
			break
		}
		if err := o.Wait(ctx); err != nil {
			log.WithError(err).Error("Failed to wait for oracle")
			return err
		}
	}
	o.Map[server.UserNameKey(userName)] = nil
	o.Unlock()

	err := func() error {
		// Check for duplicate username
		rows, err := s.GetDB().QueryContext(ctx, "SELECT id FROM user WHERE user_name = ? LIMIT 1", userName)
		if err != nil {
			log.WithError(err).Error("Failed to lookup user name")
			return server.ErrBusy
		}
		var found bool
		for rows.Next() {
			found = true
		}
		if found {
			c.SendResult(&model.RegisterResponse{
				Success: proto.Bool(false),
				Reason:  proto.String("Duplicate username"),
			})
			return nil
		}
		// Insert user
		user := &database.User{}
		user.UserName = req.UserName
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 0)
		if err != nil {
			log.WithError(err).Error("Failed to generate passowrd")
			return server.ErrBusy
		}
		user.Auth = &model.UserAuth{}
		user.Auth.PasswordHash = passwordHash
		if err := s.GetDB().InsertUser(ctx, user); err != nil {
			log.WithError(err).Error("Failed to insert user")
			return server.ErrBusy
		}
		c.SendResult(&model.RegisterResponse{
			Success: proto.Bool(true),
		})
		return nil
	}()

	// Unlock username
	o.Lock()
	delete(o.Map, server.UserNameKey(userName))
	o.Broadcast()
	o.Unlock()
	return err
}

func init() {
	router.Action("/api/register", Handle_Register)
}

package service

import (
	"context"
	"database/sql"

	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/model/common"
	kafkapb "github.com/syzoj/syzoj-ng-go/service/user/kafka"
	"github.com/syzoj/syzoj-ng-go/service/user/rpc"
	"github.com/syzoj/syzoj-ng-go/util"
)

var log = logrus.StandardLogger()

func (c *serv) RegisterUser(ctx context.Context, req *rpc.RegisterUserRequest) (*rpc.RegisterUserResponse, error) {
	userName, password := req.GetUserName(), req.GetPassword()
	if !checkUserName(userName) {
		return &rpc.RegisterUserResponse{Error: rpc.Error_InvalidUserName.Enum()}, nil
	}
	txn, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.WithError(err).Error("Failed to start transaction")
		return nil, err
	}
	var val int
	err = txn.QueryRowContext(ctx, "SELECT 1 FROM users WHERE user_name=?", userName).Scan(&val)
	if err == nil {
		return &rpc.RegisterUserResponse{Error: rpc.Error_DuplicateUserName.Enum()}, nil
	} else if err != sql.ErrNoRows {
		log.WithError(err).Error("Failed to query database")
		return nil, err
	}
	userId := util.RandomString(12)
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(password), 0); err != nil {
		log.WithError(err).Error("Failed to generate hash")
		return nil, err
	}

	if _, err = txn.ExecContext(ctx, "INSERT INTO users (id, user_name, hash) VALUES (?, ?, ?)", userId, userName, hash); err != nil {
		log.WithError(err).Error("Failed to insert into database")
		return nil, err
	}
	if err = txn.Commit(); err != nil {
		log.WithError(err).Error("Failed to commit")
		return nil, err
	}
	c.writeKafkaMessage(&kafkapb.UserEvent{UserId: proto.String(userId), Event: &kafkapb.UserEvent_Register{Register: &kafkapb.UserRegisterEvent{}}})
	return &rpc.RegisterUserResponse{UserId: proto.String(userId)}, nil
}

func (c *serv) LoginUser(ctx context.Context, req *rpc.LoginUserRequest) (*rpc.LoginUserResponse, error) {
	userName, password, deviceInfo := req.GetUserName(), req.GetPassword(), req.GetDeviceInfo()
	if !checkUserName(userName) {
		return &rpc.LoginUserResponse{Error: rpc.Error_InvalidUserName.Enum()}, nil
	}
	txn, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.WithError(err).Error("Failed to start transaction")
		return nil, err
	}
	var userId string
	var hash []byte
	err = txn.QueryRowContext(ctx, "SELECT id, hash FROM users WHERE user_name=?", userName).Scan(&userId, &hash)
	if err == sql.ErrNoRows {
		return &rpc.LoginUserResponse{Error: rpc.Error_UserNotFound.Enum()}, nil
	} else if err != nil {
		log.WithError(err).Error("Failed to query database")
		return nil, err
	}
	if bcrypt.CompareHashAndPassword(hash, []byte(password)) != nil {
		return &rpc.LoginUserResponse{Error: rpc.Error_PasswordIncorrect.Enum()}, nil
	}

	deviceToken := util.RandomString(36)
	if _, err = txn.ExecContext(ctx, "INSERT INTO devices (id, token, user_id, info) VALUES (?, ?, ?)", deviceToken[0:16], deviceToken[16:48], userId, util.ProtoSql{deviceInfo}); err != nil {
		log.WithError(err).Error("Failed to insert into database")
		return nil, err
	}
	if err = txn.Commit(); err != nil {
		log.WithError(err).Error("Failed to commit")
		return nil, err
	}
	return &rpc.LoginUserResponse{Token: proto.String(deviceToken)}, nil
}

func (c *serv) VerifyDevice(ctx context.Context, req *rpc.VerifyDeviceRequest) (*rpc.VerifyDeviceResponse, error) {
	token := req.GetToken()
	if len(token) != 48 {
		return &rpc.VerifyDeviceResponse{Error: rpc.Error_InvalidToken.Enum()}, nil
	}
	info := &common.DeviceInfo{}
	var userId string
	var realToken string
	if err := c.db.QueryRowContext(ctx, "SELECT token, user_id, info FROM devices WHERE token=?", token[0:16]).Scan(&realToken, &userId, util.ProtoSql{info}); err != nil {
		if err == sql.ErrNoRows {
			return &rpc.VerifyDeviceResponse{Error: rpc.Error_InvalidToken.Enum()}, nil
		}
		return nil, err
	}
	if realToken != token[16:48] {
		return &rpc.VerifyDeviceResponse{Error: rpc.Error_InvalidToken.Enum()}, nil
	}
	return &rpc.VerifyDeviceResponse{UserId: proto.String(userId), DeviceInfo: info}, nil
}

package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/golang/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
)

func Handle_Debug_Add_Judger(ctx context.Context) error {
	s := server.GetServer(ctx)
	c := server.GetApiContext(ctx)
	msg := &model.DebugAddJudgerRequest{}
	if err := c.ReadBody(msg); err != nil {
		return server.ErrBadRequest
	}
	j := &database.Judger{}
	var tokenBytes [32]byte
	rand.Read(tokenBytes[:])
	j.Token = proto.String(base64.URLEncoding.EncodeToString(tokenBytes[:]))
	if err := s.GetDB().InsertJudger(ctx, j); err != nil {
		log.WithError(err).Error("Failed to insert judger")
		return server.ErrBusy
	}
	c.Mutate("", "", &model.DebugAddJudgerResponse{
		JudgerId:    proto.String(string(j.GetId())),
		JudgerToken: j.Token,
	})
	return nil
}

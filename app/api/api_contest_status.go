package api

import (
	"io"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/valyala/fastjson"
)

type contestStatusContext struct {
	srv       *ApiServer
	contestId primitive.ObjectID
	userId    primitive.ObjectID

	wsConn *websocket.Conn
}

func Handle_Contest_Status(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var ct contestStatusContext
	ct.srv = c.Server()
	ct.contestId = contestId
	ct.userId = c.Session.AuthUserUid
	ct.wsConn, err = c.UpgradeWebSocket()
	if err != nil {
		panic(err)
	}
	ct.run()
	return
}

func (c *contestStatusContext) run() {
	var err error
	defer c.close()
	arena := new(fastjson.Arena)
	for {
		arena.Reset()
		select {
		case <-c.srv.ctx.Done():
			return
		case <-time.After(time.Second):
		}
		msg := arena.NewObject()
		func() {
			contest := c.srv.c.GetContestR(c.contestId)
			if contest == nil {
				return
			}
			defer contest.RUnlock()
			if contest.Running() {
				msg.Set("running", arena.NewTrue())
			} else {
				msg.Set("running", arena.NewFalse())
			}
		}()

		var w io.WriteCloser
		w, err = c.wsConn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Warning("Failed to write to WebSocket: ", err)
			return
		}
		_, err = w.Write(msg.MarshalTo(nil))
		if err != nil {
			log.Warning("Failed to write to WebSocket: ", err)
			return
		}
		err = w.Close()
	}
}

func (c *contestStatusContext) close() {
	c.wsConn.Close()
}

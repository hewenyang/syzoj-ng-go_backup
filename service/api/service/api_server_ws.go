package service

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WsApiContext struct {
	r          *http.Request
	wsConn     *websocket.Conn
	cancelFunc func()
	msg        chan<- proto.Message
}

func (s *apiService) WrapWsHandler(ctx context.Context, h func(context.Context, *WsApiContext) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.WithError(err).Error("Failed to upgrade websocket")
			return
		}
		msgChan := make(chan proto.Message, 10)
		c := &WsApiContext{
			r:      r,
			wsConn: wsConn,
			msg:    msgChan,
		}
		ctx, c.cancelFunc = context.WithCancel(ctx)
		defer c.cancelFunc()
		go func() {
			for {
				select {
				case <-ctx.Done():
					wsConn.Close()
					return
				case msg := <-msgChan:
					w, err := wsConn.NextWriter(websocket.TextMessage)
					if err != nil {
						c.cancelFunc()
						continue
					}
					if err := jsonMarshaler.Marshal(w, msg); err != nil {
						c.cancelFunc()
						continue
					}
					if err := w.Close(); err != nil {
						c.cancelFunc()
						continue
					}
				}
			}
		}()

		err = h(ctx, c)
		if err != nil {
			log.WithError(err).Error("Websocket error")
		}
	})
}

func (c *WsApiContext) Vars() map[string]string {
	return mux.Vars(c.r)
}

func (c *WsApiContext) GetCookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *WsApiContext) GetHeader(name string) string {
	return c.r.Header.Get(name)
}

func (c *WsApiContext) GetRemoteAddr() string {
	return c.r.RemoteAddr
}

func (c *WsApiContext) ReadBody(val proto.Message) error {
	messageType, r, err := c.wsConn.NextReader()
	if err != nil {
		c.cancelFunc()
		return err
	}
	if messageType != websocket.TextMessage {
		log.Error("Cannot use non-text message for websocket")
		return ErrBusy
	}
	return jsonUnmarshaler.Unmarshal(r, val)
}

func (c *WsApiContext) Send(val proto.Message) {
	c.msg <- val
}

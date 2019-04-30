package service

import (
	"context"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/service/api/model"
)

type ApiContext struct {
	r *http.Request
	w http.ResponseWriter
}

var jsonMarshaler = jsonpb.Marshaler{OrigName: true}
var jsonUnmarshaler = jsonpb.Unmarshaler{}

func (s *apiService) wrapHandler(ctx context.Context, h func(context.Context, *ApiContext) error, checkToken bool) http.Handler {
	if s.debug {
		checkToken = false
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{r: r, w: w}
		ctx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()
		if checkToken {
			token := r.Header.Get("X-CSRF-Token")
			if token != "1" {
				c.SendError(ErrCSRF)
				return
			}
		}
		err := h(ctx, c)
		if err != nil {
			c.SendError(err)
		}
	})
}

func (c *ApiContext) Vars() map[string]string {
	return mux.Vars(c.r)
}

func (c *ApiContext) GetCookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *ApiContext) SetCookie(k *http.Cookie) {
	http.SetCookie(c.w, k)
}

func (c *ApiContext) GetHeader(name string) string {
	return c.r.Header.Get(name)
}

func (c *ApiContext) GetRemoteAddr() string {
	return c.r.RemoteAddr
}

func (c *ApiContext) ReadBody(val proto.Message) error {
	return jsonUnmarshaler.Unmarshal(c.r.Body, val)
}

func (c *ApiContext) SendBody(val proto.Message) {
	if err := jsonMarshaler.Marshal(c.w, val); err != nil {
		log.WithError(err).Error("Failed to send response")
	}
}

func (c *ApiContext) SendError(err error) {
	c.SendBody(&model.Response{Error: proto.String(err.Error())})
}

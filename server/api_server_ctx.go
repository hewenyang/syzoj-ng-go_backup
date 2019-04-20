package server

import (
	"context"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/model"
)

type apiContextKey struct{}
type ApiContext struct {
	r   *http.Request
	w   http.ResponseWriter
	s   *ApiServer
	mut []*model.Mutation
}

var jsonMarshaler = jsonpb.Marshaler{OrigName: true}
var jsonUnmarshaler = jsonpb.Unmarshaler{}

func (s *ApiServer) WrapHandler(h func(context.Context) error, checkToken bool) http.Handler {
	if s.debug {
		checkToken = false
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{r: r, w: w, s: s}
		ctx := s.ctx
		ctx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()
		ctx = context.WithValue(ctx, apiContextKey{}, c)
		defer c.Send()
		if checkToken {
			token := r.Header.Get("X-CSRF-Token")
			if token != "1" {
				c.SendError(ErrCSRF)
				return
			}
		}
		err := h(ctx)
		if err != nil {
			c.SendError(err)
		}
	})
}

func (s *ApiServer) WrapDebugHandler(h func(context.Context) error) http.Handler {
	if !s.debug {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{r: r, w: w, s: s}
		ctx := s.ctx
		ctx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()
		ctx = context.WithValue(ctx, apiContextKey{}, c)
		defer c.Send()
		err := h(ctx)
		if err != nil {
			c.SendError(err)
		}
	})
}

func GetApiContext(ctx context.Context) *ApiContext {
	return ctx.Value(apiContextKey{}).(*ApiContext)
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

func (c *ApiContext) Mutate(path string, method string, val proto.Message) {
	m, err := types.MarshalAny(val)
	if err != nil {
		log.WithError(err).Error("Failed to marshal message into any")
		return
	}
	mutation := &model.Mutation{
		Path:   proto.String(path),
		Method: proto.String(method),
		Value:  m,
	}
	c.mut = append(c.mut, mutation)
}

func (c *ApiContext) SendBody(val proto.Message) {
	c.Mutate("", "setBody", val)
}

func (c *ApiContext) SendError(err error) {
	c.Mutate("", "setError", &model.Error{Error: proto.String(err.Error())})
}

func (c *ApiContext) Redirect(path string) {
	c.Mutate("", "redirect", &model.Path{Path: proto.String(path)})
}

func (c *ApiContext) Send() {
	resp := &model.Response{Mutations: c.mut}
	if err := jsonMarshaler.Marshal(c.w, resp); err != nil {
		log.WithError(err).Error("Failed to send response")
	}
}

package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/syzoj/syzoj-ng-go/model"
)

type ApiServer struct {
	s *Server

	debug       bool
	ctx         context.Context
	router      *mux.Router
	wsUpgrader  websocket.Upgrader
	wg          sync.WaitGroup
	wsConn      map[*websocket.Conn]struct{}
	wsConnMutex sync.Mutex
	cancelFunc  func()
}

type apiContextKey struct{}
type ApiContext struct {
	r   *http.Request
	w   http.ResponseWriter
	s   *ApiServer
	mut []*model.Mutation
}

var jsonMarshaler = jsonpb.Marshaler{OrigName: true}
var jsonUnmarshaler = jsonpb.Unmarshaler{}

type ApiConfig struct {
	Debug bool `json:"debug"`
}

func (s *Server) newApiServer(cfg *ApiConfig) *ApiServer {
	ApiServer := new(ApiServer)
	ApiServer.s = s
	if cfg.Debug {
		ApiServer.debug = true
	}
	router := mux.NewRouter()
	router.PathPrefix("/api").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
		w.WriteHeader(200)
	})
	ApiServer.router = router
	ApiServer.ctx, ApiServer.cancelFunc = context.WithCancel(s.ctx)
	ApiServer.wg.Add(1)
	return ApiServer
}

func (s *Server) ApiServer() *ApiServer {
	return s.apiServer
}

func (s *ApiServer) close() {
	s.cancelFunc()
	s.wsConnMutex.Lock()
	for conn := range s.wsConn {
		conn.Close()
	}
	s.wsConn = nil
	s.wsConnMutex.Unlock()
	s.wg.Done()
	s.wg.Wait()
}

func (s *ApiServer) Router() *mux.Router {
	return s.router
}

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

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.wg.Add(1)
	defer s.wg.Done()
	curTime := time.Now()
	defer func() {
		d := time.Now().Sub(curTime)
		log.Info(r)
		log.Info("Request spent ", d, int64(d))
	}()
	if s.debug {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	s.router.ServeHTTP(w, r)
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
	m, err := ptypes.MarshalAny(val)
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

func (c *ApiContext) UpgradeWebSocket() (*websocket.Conn, error) {
	return c.s.wsUpgrader.Upgrade(c.w, c.r, nil)
}

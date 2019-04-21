package handlers

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/server"
)

var log = logrus.StandardLogger()

type routeList struct {
	routes []*route
}

type route struct {
	path      string
	handler   func(context.Context) error
	websocket bool
	post      bool
	debug     bool
}

func (r *routeList) Get(path string, handler func(context.Context) error) {
	r.routes = append(r.routes, &route{path: path, handler: handler, post: false, debug: false})
}
func (r *routeList) Action(path string, handler func(context.Context) error) {
	r.routes = append(r.routes, &route{path: path, handler: handler, post: true, debug: false})
}
func (r *routeList) GetDebug(path string, handler func(context.Context) error) {
	r.routes = append(r.routes, &route{path: path, handler: handler, post: false, debug: true})
}
func (r *routeList) ActionDebug(path string, handler func(context.Context) error) {
	r.routes = append(r.routes, &route{path: path, handler: handler, post: true, debug: true})
}
func (r *routeList) WebSocket(path string, handler func(context.Context) error) {
	r.routes = append(r.routes, &route{path: path, handler: handler, websocket: true})
}

var router = &routeList{}

func RegisterHandlers(s *server.ApiServer) {
	r := s.Router()
	for _, route := range router.routes {
		m := r.Path(route.path)
		if route.post {
			m = m.Methods("POST")
		} else {
			m = m.Methods("GET")
		}
		if route.websocket {
			m = m.Handler(s.WrapWsHandler(route.handler))
		} else if route.debug {
			m = m.Handler(s.WrapDebugHandler(route.handler))
		} else {
			m = m.Handler(s.WrapHandler(route.handler, true))
		}
	}
	r.PathPrefix("/api").Methods("GET").Handler(s.WrapHandler(Handle_Not_Found, true))
}

func Handle_Not_Found(ctx context.Context) error {
	return server.ErrNotFound
}

package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ApiServer struct {
	s *Server

	debug      bool
	ctx        context.Context
	router     *mux.Router
	wsUpgrader websocket.Upgrader
	wg         sync.WaitGroup
	cancelFunc func()
}

type ApiConfig struct {
	Debug bool `json:"debug"`
}

func (s *Server) newApiServer(cfg *ApiConfig) *ApiServer {
	server := &ApiServer{}
	server.s = s
	if cfg.Debug {
		server.debug = true
		server.wsUpgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	router := mux.NewRouter()
	router.PathPrefix("/api").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
		w.WriteHeader(200)
	})
	server.router = router
	server.ctx, server.cancelFunc = context.WithCancel(s.WithServer(context.Background()))
	server.wg.Add(1)
	return server
}

func (s *Server) ApiServer() *ApiServer {
	return s.apiServer
}

func (s *ApiServer) close() {
	s.cancelFunc()
	s.wg.Done()
	s.wg.Wait()
}

func (s *ApiServer) Router() *mux.Router {
	return s.router
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

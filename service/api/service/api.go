package service

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (s *apiService) setupApiRouter(ctx context.Context, router *mux.Router) {
	s.setupPost(ctx, router, "/login", s.handleLogin)
	s.setupPost(ctx, router, "/register", s.handleRegister)
	s.setupGet(ctx, router, "/problems", s.handleProblems)
}

func (s *apiService) setupGet(ctx context.Context, router *mux.Router, path string, f func(context.Context, *ApiContext) error) {
	router.Path(path).Handler(s.wrapHandler(ctx, f, true)).Methods("GET")
}

func (s *apiService) setupPost(ctx context.Context, router *mux.Router, path string, f func(context.Context, *ApiContext) error) {
	router.Path(path).Handler(s.wrapHandler(ctx, f, true)).Methods("POST")
}

func (s *apiService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.wg.Add(1)
	defer s.wg.Done()
	curTime := time.Now()
	defer func() {
		d := time.Now().Sub(curTime)
		log.Info(r)
		log.Info("Request spent ", d)
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

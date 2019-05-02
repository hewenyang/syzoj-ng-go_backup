package service

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/service"
	problemclient "github.com/syzoj/syzoj-ng-go/service/problem/client"
	userclient "github.com/syzoj/syzoj-ng-go/service/user/client"
)

type Config struct {
	Debug    bool
	HttpAddr string
}

type apiService struct {
	config        *Config
	log           *logrus.Logger
	userClient    *userclient.Client
	problemClient *problemclient.Client
	addr          string
	debug         bool
	wsUpgrader    websocket.Upgrader
	router        *mux.Router
	wg            sync.WaitGroup
}

func NewApiService(config *Config) *service.ServiceInfo {
	return service.ServiceBuilder{
		Name:    "api",
		Version: "v0.0.1",
		Object:  &apiService{config: config},
	}.Build()
}

func (s *apiService) Main(ctx context.Context, c *service.ServiceContext) {
	var err error
	s.log = c.GetLogger()
	if s.userClient, err = userclient.NewUserClient(); err != nil {
		s.log.WithError(err).Error("Failed to connect to user service")
		return
	}
	defer s.userClient.Close()
	if s.problemClient, err = problemclient.NewProblemClient(); err != nil {
		s.log.WithError(err).Error("Failed to connect to problem service")
		return
	}
	defer s.problemClient.Close()
	s.debug = s.config.Debug
	s.addr = s.config.HttpAddr
	s.config = nil
	if s.debug {
		s.wsUpgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	s.router = mux.NewRouter()
	s.router.Handle("/", http.FileServer(http.Dir("static")))
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	s.setupApiRouter(ctx, apiRouter)

	server := &http.Server{
		Handler:      s,
		WriteTimeout: time.Second * 10,
	}
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.log.WithError(err).Error("Failed to listen")
		return
	}
	s.wg.Add(1)
	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			s.log.WithError(err).Error("Failed to serve HTTP")
		}
		s.wg.Done()
	}()
	go func() {
		<-ctx.Done()
		server.Close()
	}()
	c.StartupDone()
	s.wg.Wait()
}

func (s *apiService) Migrate(ctx context.Context, c *service.ServiceContext, prevVersion string) error {
	s.log = c.GetLogger()
	s.log.Infof("apiService: Migrating from %s to %s", prevVersion, "v0.0.1")
	return nil
}

package server

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/database"
)

var log = logrus.StandardLogger()

type Server struct {
	db     *database.Database
	oracle *Oracle

	apiServer   *ApiServer
	judgeServer *JudgeServer

	runningJudges sync.Map // database.Submission -> *JudgeRequest
}

type serverKey struct{}

type ServerConfig struct {
	API ApiConfig `json:"api"`
}

func NewServer(db *database.Database, cfg *ServerConfig) *Server {
	server := new(Server)
	server.db = db
	server.oracle = NewOracle()
	server.apiServer = server.newApiServer(&cfg.API)
	server.judgeServer = server.newJudgeServer()
	return server
}

func (s *Server) WithServer(ctx context.Context) context.Context {
	return context.WithValue(ctx, serverKey{}, s)
}

func GetServer(ctx context.Context) *Server {
	return ctx.Value(serverKey{}).(*Server)
}

func (s *Server) GetDB() *database.Database {
	return s.db
}

func (s *Server) GetOracle() *Oracle {
	return s.oracle
}

func (s *Server) GetJudge() *JudgeServer {
	return s.judgeServer
}

func (s *Server) GetContext() context.Context {
	return context.TODO()
}

func (s *Server) Close() error {
	s.apiServer.close()
	s.judgeServer.close()
	return nil
}

func (s *Server) GetRunningJudges() *sync.Map {
	return &s.runningJudges
}

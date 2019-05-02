package service

import (
	"container/list"
	"context"
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/syzoj/syzoj-ng-go/fakenet"
	"github.com/syzoj/syzoj-ng-go/service"
	"github.com/syzoj/syzoj-ng-go/service/judge/rpc"
	"github.com/syzoj/syzoj-ng-go/util"
)

var log = logrus.StandardLogger()

type Config struct {
	MySQL string
}

func NewJudgeService(config *Config) *service.ServiceInfo {
	return service.ServiceBuilder{
		Name:    "judge",
		Version: "v0.0.1",
		Object:  &serv{config: config},
	}.Build()
}

type serv struct {
	config *Config

	wg   sync.WaitGroup
	db   *sql.DB
	mu   sync.Mutex
	subs map[string]*list.List
}

func (s *serv) Main(ctx context.Context, c *service.ServiceContext) {
	lis, err := fakenet.Base.Listen("service-judge")
	if err != nil {
		log.WithError(err).Error("Failed to listen")
		return
	}
	gserver := grpc.NewServer()
	rpc.RegisterJudgeServer(gserver, s)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		gserver.Serve(lis)
	}()
	c.StartupDone()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		<-ctx.Done()
		gserver.Stop()
	}()
	s.wg.Wait()
}

func (s *serv) Migrate(prevVersion string) error {
	ctx := context.Background()
	switch prevVersion {
	case "":
		return util.MigrateMySQL(ctx, s.config.MySQL, "CREATE TABLE submissions (id VARCHAR(64) PRIMARY KEY, test_data BLOB, submit_data BLOB, result BLOB);")
	}
	return nil
}

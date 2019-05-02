package service

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/syzoj/syzoj-ng-go/fakenet"
	"github.com/syzoj/syzoj-ng-go/service"
	"github.com/syzoj/syzoj-ng-go/service/judge/rpc"
	"github.com/syzoj/syzoj-ng-go/util"
)

type Config struct {
	MySQL        string
	KafkaBrokers []string
}

func NewJudgeService(config *Config) *service.ServiceInfo {
	return service.ServiceBuilder{
		Name:    "judge",
		Version: "v0.0.1",
		Object:  &serv{config: config},
	}.Build()
}

type serv struct {
	config      *Config
	log         *logrus.Logger
	wg          sync.WaitGroup
	ctx         context.Context
	db          *sql.DB
	kafkaWriter *kafka.Writer
	mu          sync.Mutex
	subs        map[string]*subEntry
}

func (s *serv) Main(ctx context.Context, c *service.ServiceContext) {
	var err error
	s.log = c.GetLogger()
	s.db, err = sql.Open("mysql", s.config.MySQL)
	if err != nil {
		s.log.WithError(err).Error("Failed to connect to MySQL")
		return
	}
	defer s.db.Close()
	s.kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  s.config.KafkaBrokers,
		Topic:    "judge",
		Balancer: &kafka.Hash{},
	})
	defer s.kafkaWriter.Close()
	lis, err := fakenet.Base.Listen("service-judge")
	if err != nil {
		s.log.WithError(err).Error("Failed to listen")
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

func (s *serv) Migrate(ctx context.Context, c *service.ServiceContext, prevVersion string) error {
	s.log = c.GetLogger()
	switch prevVersion {
	case "":
		return util.MigrateMySQL(ctx, s.config.MySQL, "CREATE TABLE submissions (id VARCHAR(64) PRIMARY KEY, test_data BLOB, submit_data BLOB, result BLOB, done BOOLEAN);")
	}
	return nil
}

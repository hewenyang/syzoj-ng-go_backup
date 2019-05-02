package service

import (
	"context"
	"database/sql"
	"net"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/syzoj/syzoj-ng-go/fakenet"
	"github.com/syzoj/syzoj-ng-go/service"
	judgecli "github.com/syzoj/syzoj-ng-go/service/judge/client"
	"github.com/syzoj/syzoj-ng-go/service/problem/rpc"
)

type Config struct {
	MySQL        string
	KafkaBrokers []string
}

type serv struct {
	config      *Config
	ctx         context.Context
	log         *logrus.Logger
	wg          sync.WaitGroup
	db          *sql.DB
	kafkaWriter *kafka.Writer
	judgeCli    *judgecli.Client
}

func NewProblemService(config *Config) *service.ServiceInfo {
	return service.ServiceBuilder{
		Name:    "problem",
		Version: "v0.0.1",
		Object:  &serv{config: config},
	}.Build()
}

func (s *serv) Main(ctx context.Context, c *service.ServiceContext) {
	var err error
	s.ctx = ctx
	s.log = c.GetLogger()
	if s.db, err = sql.Open("mysql", s.config.MySQL); err != nil {
		s.log.WithError(err).Error("Failed to open MySQL")
		return
	}
	defer s.db.Close()
	s.kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  s.config.KafkaBrokers,
		Topic:    "problem",
		Balancer: &kafka.Hash{},
	})
	defer s.kafkaWriter.Close()
	if s.judgeCli, err = judgecli.NewJudgeClient(); err != nil {
		s.log.WithError(err).Error("Failed to connect to judge service")
		return
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterProblemServer(grpcServer, s)
	var listener net.Listener
	if listener, err = fakenet.Base.Listen("service-problem"); err != nil {
		s.log.WithError(err).Error("Failed to listen")
		return
	}
	c.StartupDone()
	s.wg.Add(1)
	go func() {
		<-ctx.Done()
		grpcServer.Stop()
		s.wg.Done()
	}()
	grpcServer.Serve(listener)
	s.wg.Wait()
}

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
	"github.com/syzoj/syzoj-ng-go/service/user/rpc"
)

type Config struct {
	MySQL        string
	KafkaBrokers []string
}

type serv struct {
	config      *Config
	log         *logrus.Logger
	ctx         context.Context
	wg          sync.WaitGroup
	db          *sql.DB
	kafkaWriter *kafka.Writer
}

func NewUserService(config *Config) *service.ServiceInfo {
	return service.ServiceBuilder{
		Name:    "user",
		Version: "v0.0.1",
		Object:  &serv{config: config},
	}.Build()
}

func (s *serv) Main(ctx context.Context, c *service.ServiceContext) {
	var err error
	s.log = c.GetLogger()
	s.ctx = ctx
	if s.db, err = sql.Open("mysql", s.config.MySQL); err != nil {
		s.log.WithError(err).Error("Failed to open MySQL")
		return
	}
	defer s.db.Close()
	s.kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  s.config.KafkaBrokers,
		Topic:    "user",
		Balancer: &kafka.Hash{},
	})
	defer s.kafkaWriter.Close()
	grpcServer := grpc.NewServer()
	rpc.RegisterUserServer(grpcServer, s)
	var listener net.Listener
	if listener, err = fakenet.Base.Listen("service-user"); err != nil {
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

package service

import (
	"context"
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/proto"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/syzoj/syzoj-ng-go/fakenet"
	"github.com/syzoj/syzoj-ng-go/service"
	"github.com/syzoj/syzoj-ng-go/service/user-profile/rpc"
	userkpb "github.com/syzoj/syzoj-ng-go/service/user/kafka"
)

var log = logrus.StandardLogger()

type Config struct {
	MySQL        string
	KafkaBrokers []string
}

type serv struct {
	config      *Config
	wg          sync.WaitGroup
	db          *sql.DB
	kafkaReader *kafka.Reader
}

func NewUserProfileService(config *Config) *service.ServiceInfo {
	return service.ServiceBuilder{
		Name:    "user-profile",
		Version: "v0.0.1",
		Object:  &serv{config: config},
	}.Build()
}

func (s *serv) Main(ctx context.Context, c *service.ServiceContext) {
	var err error
	s.db, err = sql.Open("mysql", s.config.MySQL)
	if err != nil {
		log.WithError(err).Error("Failed to connect to MySQL")
		return
	}
	defer s.db.Close()
	s.kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        s.config.KafkaBrokers,
		Topic:          "user",
		GroupID:        "user-profile",
		CommitInterval: time.Second,
	})
	defer s.kafkaReader.Close()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			m, err := s.kafkaReader.FetchMessage(ctx)
			if err != nil {
				break
			}
			msg := &userkpb.UserEvent{}
			if err := proto.Unmarshal(m.Value, msg); err != nil {
				log.WithError(err).Error("Failed to unmarshal message from topic:user to UserEvent")
				continue
			}
		retry:
			if err := s.handleUserEvent(ctx, msg); err != nil {
				log.WithError(err).Error("Failed to process message")
				t := time.NewTimer(time.Second * 1)
				select {
				case <-ctx.Done():
					t.Stop()
					break
				case <-t.C:
					goto retry
				}
			}
			if err := s.kafkaReader.CommitMessages(ctx, m); err != nil {
				log.WithError(err).Error("Failed to commit message")
			}
		}
	}()
	listener, err := fakenet.Base.Listen("service-user-profile")
	if err != nil {
		log.WithError(err).Error("Failed to listen on service-user-profile")
		return
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterUserProfileServer(grpcServer, s)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := grpcServer.Serve(listener); err != nil {
			log.WithError(err).Error("Failed to serve gRPC")
		}
	}()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		<-ctx.Done()
		grpcServer.Stop()
	}()
	c.StartupDone()
	s.wg.Wait()
}

package service

import (
	"github.com/gogo/protobuf/proto"
	"github.com/segmentio/kafka-go"

	kafkapb "github.com/syzoj/syzoj-ng-go/service/user/kafka"
)

func (s *serv) writeKafkaMessage(msg *kafkapb.UserEvent) {
	b, err := proto.Marshal(msg)
	if err != nil {
		log.WithError(err).Error("Failed to marshal protobuf message")
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.kafkaWriter.WriteMessages(s.ctx, kafka.Message{
			Key:   []byte(msg.GetUserId()),
			Value: b,
		}); err != nil {
			log.WithError(err).Error("Failed to write user kafka message")
		}
	}()
}

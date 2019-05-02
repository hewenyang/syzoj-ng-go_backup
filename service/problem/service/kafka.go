package service

import (
	"github.com/gogo/protobuf/proto"
	"github.com/segmentio/kafka-go"

	kafkapb "github.com/syzoj/syzoj-ng-go/service/problem/kafka"
)

func (s *serv) writeKafkaMessage(msg *kafkapb.ProblemEvent) {
	b, err := proto.Marshal(msg)
	if err != nil {
		s.log.WithError(err).Error("Failed to marshal protobuf message")
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.kafkaWriter.WriteMessages(s.ctx, kafka.Message{
			Key:   []byte(msg.GetProblemId()),
			Value: b,
		}); err != nil {
			s.log.WithError(err).Error("Failed to write kafka message")
		}
	}()
}

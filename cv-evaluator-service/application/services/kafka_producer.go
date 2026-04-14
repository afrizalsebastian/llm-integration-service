package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/afrizalsebastian/go-common-modules/kafka"
	"github.com/afrizalsebastian/go-common-modules/logger"
)

type IKafkaProducer interface {
	PublishMessage(ctx context.Context, topic string, key, message interface{}) error
}

type kafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(producer *kafka.Producer) IKafkaProducer {
	return &kafkaProducer{
		producer: producer,
	}
}

func (k *kafkaProducer) PublishMessage(ctx context.Context, topic string, key, message interface{}) error {
	l := logger.New().WithContext(ctx)

	if k.producer == nil {
		l.Warn("kafka not initialized").Msg()
		return errors.New("kafka not initialized")
	}

	byteKey, _ := json.Marshal(key)
	byteMessage, _ := json.Marshal(message)

	_, _, err := k.producer.Publish(topic, byteKey, byteMessage)

	if err != nil {
		l.Warn("kafka failed to publish message").Msg()
		return err
	}

	return nil
}

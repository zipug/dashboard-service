package kafka

import (
	"dashboard/internal/common/service/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaRepository struct {
	producer *kafka.Producer
}

func NewKafkaRepository(cfg *config.AppConfig) *KafkaRepository {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Broker,
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil
	}
	return &KafkaRepository{producer: p}
}

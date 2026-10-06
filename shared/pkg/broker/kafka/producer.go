package kafka

import (
	"context"
	"errors"
	"sync"

	"github.com/egor200512/URL_shortener/shared/configs"
)

type Producer struct {
	client    client
	topic     string
	closeOnce sync.Once
}

func NewProducer(conf *configs.KafkaConf) (*Producer, error) {
	kafkaClient, err := newProducerClient(conf)
	if err != nil {
		return nil, err
	}
	return &Producer{client: kafkaClient, topic: conf.Topic()}, nil
}

func (p *Producer) Publish(ctx context.Context, topic, key string, payload []byte) error {
	if p == nil || p.client == nil {
		return errors.New("kafka producer is not initialized")
	}
	if topic == "" {
		return errors.New("kafka topic is required")
	}
	if key == "" {
		return errors.New("kafka message key is required")
	}

	return p.client.Produce(ctx, &record{
		topic:     topic,
		partition: -1,
		key:       []byte(key),
		value:     append([]byte(nil), payload...),
	})
}

func (p *Producer) Topic() string {
	if p == nil {
		return ""
	}
	return p.topic
}

func (p *Producer) Close() {
	if p == nil || p.client == nil {
		return
	}
	p.closeOnce.Do(p.client.Close)
}

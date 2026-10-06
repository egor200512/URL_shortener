package broker

import "context"

type IProducer interface {
	Publish(ctx context.Context, topic, key string, payload []byte) error
	Topic() string
	Close()
}

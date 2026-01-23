package broker

import (
	"context"
)

type IProducer interface {
	Publish(ctx context.Context, subject string, data []byte) error
	Close()
	Prefix() string
}

type IConsumer interface {
	Read()
	Close()
	Prefix() string
}

package broker

import (
	"context"
	"time"

	gonats "github.com/nats-io/nats.go"
)

type IProducer interface {
	Publish(ctx context.Context, subject string, data []byte) error
	Close()
	Prefix() string
}

type IConsumer interface {
	Fetch(batch int, wait time.Duration) ([]*gonats.Msg, error)
	Close()
}

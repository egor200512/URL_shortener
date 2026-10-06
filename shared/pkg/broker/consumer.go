package broker

import "context"

type MessageHandler func(ctx context.Context, payload []byte) error

type IConsumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
	Close()
}

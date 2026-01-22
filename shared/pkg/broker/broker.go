package broker

import (
	"context"
)

type IBroker interface {
	Publish(ctx context.Context, subject string, data []byte) error
	Close()
	Prefix() string
}

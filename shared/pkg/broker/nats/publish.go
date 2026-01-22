package nats

import (
	"context"
)

func (b *NatsBroker) Publish(ctx context.Context, subject string, data []byte) error {
	conn, err := b.connOrErr()
	if err != nil {
		return err
	}

	return conn.Publish(subject, data)
}

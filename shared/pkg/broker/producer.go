package broker

import "context"

type IProducer interface {
	Publish(ctx context.Context, subject string, payload []byte) error
	Subject(eventType string) string
	Close()
}

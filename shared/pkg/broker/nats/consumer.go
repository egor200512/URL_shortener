package nats

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
	"github.com/egor200512/URL_shortener/shared/pkg/logger"
	"github.com/nats-io/nats.go"
)

type Consumer struct {
	conn          *nats.Conn
	js            nats.JetStreamContext
	subscription  *nats.Subscription
	stream        string
	subjectPrefix string
	consumerName  string
	batch         int
	wait          time.Duration
}

func NewConsumer(conf *configs.NatsConf) *Consumer {
	conn, err := nats.Connect(conf.URL())
	if err != nil {
		logger.Fatal("failed to connect to nats", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		logger.Fatal("failed to create nats jetstream context", err)
	}

	c := &Consumer{
		conn:          conn,
		js:            js,
		stream:        conf.Stream(),
		subjectPrefix: conf.SubjectPrefix(),
		consumerName:  conf.ConsumerName(),
		batch:         conf.ConsumerBatch(),
		wait:          conf.ConsumerWait(),
	}

	if err := c.ensureStream(); err != nil {
		logger.Fatal("failed to ensure nats stream", err)
	}

	sub, err := c.js.PullSubscribe(c.Subject(">"), c.consumerName, nats.BindStream(c.stream))
	if err != nil {
		logger.Fatal("failed to subscribe to nats stream", err)
	}
	c.subscription = sub

	return c
}

func (c *Consumer) Consume(ctx context.Context, handler broker.MessageHandler) error {
	if handler == nil {
		return errors.New("nats message handler is not initialized")
	}

	if c == nil || c.subscription == nil {
		return errors.New("nats consumer is not initialized")
	}

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		messages, err := c.subscription.Fetch(c.batch, nats.MaxWait(c.wait))
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				continue
			}
			return err
		}

		c.handleMessages(ctx, messages, handler)
	}
}

func (c *Consumer) handleMessages(ctx context.Context, messages []*nats.Msg, handler broker.MessageHandler) {
	for _, message := range messages {
		if err := handler(ctx, message.Data); err != nil {
			slog.Warn("failed to handle nats message", "error", err)
			if nakErr := message.Nak(); nakErr != nil {
				slog.Warn("failed to nak nats message", "error", nakErr)
			}
			continue
		}

		if err := message.Ack(); err != nil {
			slog.Warn("failed to ack nats message", "error", err)
		}
	}
}

func (c *Consumer) Subject(eventType string) string {
	return fmt.Sprintf("%s.%s", c.subjectPrefix, eventType)
}

func (c *Consumer) Close() {
	if c == nil {
		return
	}
	if c.subscription != nil {
		if err := c.subscription.Drain(); err != nil {
			slog.Warn("failed to drain nats subscription", "error", err)
		}
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) ensureStream() error {
	if _, err := c.js.StreamInfo(c.stream); err == nil {
		return nil
	} else if !errors.Is(err, nats.ErrStreamNotFound) {
		return err
	}

	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:     c.stream,
		Subjects: []string{c.Subject(">")},
		Storage:  nats.FileStorage,
	})
	return err
}

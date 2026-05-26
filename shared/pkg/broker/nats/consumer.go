package nats

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
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
		log.Fatalf("failed to connect to nats: %s", err.Error())
	}

	js, err := conn.JetStream()
	if err != nil {
		log.Fatalf("failed to create nats jetstream context: %s", err.Error())
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
		log.Fatalf("failed to ensure nats stream: %s", err.Error())
	}

	sub, err := c.js.PullSubscribe(c.Subject(">"), c.consumerName, nats.BindStream(c.stream))
	if err != nil {
		log.Fatalf("failed to subscribe to nats stream: %s", err.Error())
	}
	c.subscription = sub

	return c
}

func (c *Consumer) Consume(ctx context.Context, handler broker.MessageHandler) error {
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

		for _, message := range messages {
			if err := handler(ctx, message.Data); err != nil {
				log.Printf("failed to handle nats message: %s\n", err.Error())
				if nakErr := message.Nak(); nakErr != nil {
					log.Printf("failed to nak nats message: %s\n", nakErr.Error())
				}
				continue
			}

			if err := message.Ack(); err != nil {
				log.Printf("failed to ack nats message: %s\n", err.Error())
			}
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
			log.Printf("failed to drain nats subscription: %s\n", err.Error())
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

package nats

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/nats-io/nats.go"
)

type Producer struct {
	conn          *nats.Conn
	js            nats.JetStreamContext
	stream        string
	subjectPrefix string
}

func NewProducer(conf *configs.NatsConf) *Producer {
	conn, err := nats.Connect(conf.URL())
	if err != nil {
		log.Fatalf("failed to connect to nats: %s", err.Error())
	}

	js, err := conn.JetStream()
	if err != nil {
		log.Fatalf("failed to create nats jetstream context: %s", err.Error())
	}

	p := &Producer{
		conn:          conn,
		js:            js,
		stream:        conf.Stream(),
		subjectPrefix: conf.SubjectPrefix(),
	}

	if err := p.ensureStream(); err != nil {
		log.Fatalf("failed to ensure nats stream: %s", err.Error())
	}

	return p
}

func (p *Producer) Publish(ctx context.Context, subject string, payload []byte) error {
	if p == nil || p.js == nil {
		return errors.New("nats producer is not initialized")
	}

	_, err := p.js.Publish(subject, payload, nats.Context(ctx))
	return err
}

func (p *Producer) Subject(eventType string) string {
	return fmt.Sprintf("%s.%s", p.subjectPrefix, eventType)
}

func (p *Producer) Close() {
	if p == nil || p.conn == nil {
		return
	}
	p.conn.Close()
}

func (p *Producer) ensureStream() error {
	if _, err := p.js.StreamInfo(p.stream); err == nil {
		return nil
	} else if !errors.Is(err, nats.ErrStreamNotFound) {
		return err
	}

	_, err := p.js.AddStream(&nats.StreamConfig{
		Name:     p.stream,
		Subjects: []string{p.Subject(">")},
		Storage:  nats.FileStorage,
	})
	return err
}

package nats

import (
	"errors"
	"fmt"
	"log"

	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
	gonats "github.com/nats-io/nats.go"
)

type NatsBroker struct {
	conn    *gonats.Conn
	subject string
}

func NewNatsProducer(conf *configs.NatsConf) broker.IProducer {
	conn, err := gonats.Connect(conf.URL())
	if err != nil {
		log.Fatalf("failed to create nats connection: %s", err.Error())
	}

	js, err := conn.JetStream()
	if err != nil {
		log.Fatalf("failed to init jetstream: %s", err.Error())
	}

	if err := ensureLinksStream(js, conf.SubjectPrefix()); err != nil {
		log.Fatalf("failed to ensure links stream: %s", err.Error())
	}

	return &NatsBroker{
		conn:    conn,
		subject: conf.SubjectPrefix(),
	}

}

func (b *NatsBroker) connOrErr() (*gonats.Conn, error) {
	if b == nil {
		return nil, errors.New("nats broker is nil")
	}
	if b.conn == nil {
		return nil, errors.New("nats broker connection is nil")
	}
	return b.conn, nil
}

func (b *NatsBroker) Prefix() string {
	return b.subject
}

func ensureLinksStream(js gonats.JetStreamContext, prefix string) error {
	_, err := js.StreamInfo(prefix)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gonats.ErrStreamNotFound) {
		return err
	}

	subject := fmt.Sprintf("%s.*", prefix)

	_, err = js.AddStream(&gonats.StreamConfig{
		Name:        prefix,
		Subjects:    []string{subject},
		Description: "stores link events",
		Storage:     gonats.FileStorage,
		Replicas:    1,
	})

	return err
}

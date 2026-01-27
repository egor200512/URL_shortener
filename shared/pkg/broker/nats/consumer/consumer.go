package consumer

import (
	"fmt"
	"log"

	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
	gonats "github.com/nats-io/nats.go"
)

type consumer struct {
	conn *gonats.Conn
	js   gonats.JetStreamContext
	sub  *gonats.Subscription
}

func NewConsumer(conf *configs.NatsConf) broker.IConsumer {
	conn, err := gonats.Connect(conf.URL())
	if err != nil {
		log.Fatalf("failed to create nats connection: %s", err.Error())
	}

	js, err := conn.JetStream()
	if err != nil {
		log.Fatalf("failed to init jetstream: %s", err.Error())
	}

	subject := fmt.Sprintf("%s.*", conf.SubjectPrefix())

	sub, err := js.PullSubscribe(
		subject,
		conf.ConsumerName(),
		gonats.BindStream(conf.SubjectPrefix()),
		gonats.ManualAck(),
	)
	if err != nil {
		conn.Close()
		log.Fatalf("pull subscribe: %v", err)
	}

	return &consumer{conn: conn, js: js, sub: sub}
}

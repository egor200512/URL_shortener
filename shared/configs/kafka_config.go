package configs

import (
	"errors"
	"net"
	"os"
)

const (
	KAFKA_HOST           = "KAFKA_HOST"
	KAFKA_INTERNAL_PORT  = "KAFKA_INTERNAL_PORT"
	KAFKA_EXTERNAL_PORT  = "KAFKA_EXTERNAL_PORT"
)

type KafkaConf struct {
	host         string
	internalPort string
	externalPort string
}

func NewKafkaConf() (*KafkaConf, error) {
	var host, internalPort, externalPort string

	if host = os.Getenv(KAFKA_HOST); len(host) == 0 {
		return nil, errors.New("failed to get kafka host")
	}

	if internalPort = os.Getenv(KAFKA_INTERNAL_PORT); len(internalPort) == 0 {
		return nil, errors.New("failed to get kafka internal port")
	}

	if externalPort = os.Getenv(KAFKA_EXTERNAL_PORT); len(externalPort) == 0 {
		return nil, errors.New("failed to get kafka external port")
	}

	return &KafkaConf{
		host:         host,
		internalPort: internalPort,
		externalPort: externalPort,
	}, nil
}

func (conf *KafkaConf) InternalBroker() string {
	return net.JoinHostPort(conf.host, conf.internalPort)
}

func (conf *KafkaConf) ExternalBroker() string {
	return net.JoinHostPort(conf.host, conf.externalPort)
}

func (conf *KafkaConf) Brokers() []string {
	return []string{conf.InternalBroker()}
}

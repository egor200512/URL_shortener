package configs

import (
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	NATS_URL            = "NATS_URL"
	NATS_HOST           = "NATS_HOST"
	NATS_PORT           = "NATS_PORT"
	NATS_USER           = "NATS_USER"
	NATS_PASSWORD       = "NATS_PASSWORD"
	NATS_STREAM         = "NATS_STREAM"
	NATS_SUBJECT_PREFIX = "NATS_SUBJECT_PREFIX"
)

type NatsConf struct {
	url           string
	stream        string
	subjectPrefix string
}

func NewNatsConf() (*NatsConf, error) {
	url := os.Getenv(NATS_URL)
	if url == "" {
		host := os.Getenv(NATS_HOST)
		if host == "" {
			return nil, errors.New("failed to get nats host")
		}

		port := os.Getenv(NATS_PORT)
		if port == "" {
			return nil, errors.New("failed to get nats port")
		}

		user := os.Getenv(NATS_USER)
		password := os.Getenv(NATS_PASSWORD)
		credentials := ""
		if user != "" || password != "" {
			credentials = fmt.Sprintf("%s:%s@", user, password)
		}

		url = fmt.Sprintf("nats://%s%s", credentials, net.JoinHostPort(host, port))
	}

	stream := os.Getenv(NATS_STREAM)
	if stream == "" {
		return nil, errors.New("failed to get nats stream")
	}

	subjectPrefix := os.Getenv(NATS_SUBJECT_PREFIX)
	if subjectPrefix == "" {
		return nil, errors.New("failed to get nats subject prefix")
	}

	return &NatsConf{
		url:           url,
		stream:        stream,
		subjectPrefix: subjectPrefix,
	}, nil
}

func (conf *NatsConf) URL() string {
	return conf.url
}

func (conf *NatsConf) Stream() string {
	return conf.stream
}

func (conf *NatsConf) SubjectPrefix() string {
	return conf.subjectPrefix
}

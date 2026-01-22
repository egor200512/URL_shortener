package configs

import (
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	NATS_URL      = "NATS_URL"
	NATS_HOST     = "NATS_HOST"
	NATS_PORT     = "NATS_PORT"
	NATS_USER     = "NATS_USER"
	NATS_PASSWORD = "NATS_PASSWORD"
	NATS_PREFIX   = "NATS_PREFIX"
)

type NatsConf struct {
	url    string
	prefix string
}

func NewNatsConf() (*NatsConf, error) {
	var url, host, port, user, password, prefix string

	if url = os.Getenv(NATS_URL); len(url) == 0 {
		if host = os.Getenv(NATS_HOST); len(host) == 0 {
			return nil, errors.New("failed to get nats host")
		}

		if port = os.Getenv(NATS_PORT); len(port) == 0 {
			return nil, errors.New("failed to get nats port")
		}

		if user = os.Getenv(NATS_USER); len(user) == 0 {
			return nil, errors.New("failed to get nats user")
		}

		if password = os.Getenv(NATS_PASSWORD); len(password) == 0 {
			return nil, errors.New("failed to get nats password")
		}

		cred := fmt.Sprintf("%s:%s@", user, password)

		url = fmt.Sprintf("nats://%s%s", cred, net.JoinHostPort(host, port))
	}

	prefix = os.Getenv(NATS_PREFIX)

	return &NatsConf{
		url:    url,
		prefix: prefix,
	}, nil
}

func (conf *NatsConf) URL() string {
	return conf.url
}

func (conf *NatsConf) SubjectPrefix() string {
	return conf.prefix
}

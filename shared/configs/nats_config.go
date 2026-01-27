package configs

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	NATS_URL              = "NATS_URL"
	NATS_HOST             = "NATS_HOST"
	NATS_PORT             = "NATS_PORT"
	NATS_USER             = "NATS_USER"
	NATS_PASSWORD         = "NATS_PASSWORD"
	NATS_PREFIX           = "NATS_PREFIX"
	NATS_CONSUMER_NAME    = "NATS_CONSUMER_NAME"
	NATS_CONSUMER_BATCH   = "NATS_CONSUMER_BATCH"
	NATS_CONSUMER_WAIT_MS = "NATS_CONSUMER_WAIT_MS"
)

type NatsConf struct {
	url            string
	prefix         string
	consumerName   string
	consumerBatch  int
	consumerWaitMs int
}

func NewNatsConf() (*NatsConf, error) {
	var url, host, port, user, password, prefix, consumer string
	var batch, waitMs int

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

	if prefix = os.Getenv(NATS_PREFIX); len(prefix) == 0 {
		return nil, errors.New("failed to get nats prefix")
	}

	if consumer = os.Getenv(NATS_CONSUMER_NAME); len(consumer) == 0 {
		return nil, errors.New("failed to get nats consumer name")
	}

	batchStr := os.Getenv(NATS_CONSUMER_BATCH)
	if len(batchStr) == 0 {
		return nil, errors.New("failed to get nats consumer batch")
	}
	parsedBatch, err := strconv.Atoi(batchStr)
	if err != nil {
		return nil, errors.New("failed to parse nats consumer batch")
	}
	batch = parsedBatch

	waitStr := os.Getenv(NATS_CONSUMER_WAIT_MS)
	if len(waitStr) == 0 {
		return nil, errors.New("failed to get nats consumer wait")
	}
	parsedWait, err := strconv.Atoi(waitStr)
	if err != nil {
		return nil, errors.New("failed to parse nats consumer wait")
	}
	waitMs = parsedWait

	return &NatsConf{
		url:            url,
		prefix:         prefix,
		consumerName:   consumer,
		consumerBatch:  batch,
		consumerWaitMs: waitMs,
	}, nil
}

func (conf *NatsConf) URL() string {
	return conf.url
}

func (conf *NatsConf) SubjectPrefix() string {
	return conf.prefix
}

func (conf *NatsConf) ConsumerName() string {
	return conf.consumerName
}

func (conf *NatsConf) ConsumerBatch() int {
	return conf.consumerBatch
}

func (conf *NatsConf) ConsumerWait() time.Duration {
	return time.Duration(conf.consumerWaitMs) * time.Millisecond
}

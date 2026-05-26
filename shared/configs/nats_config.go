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
	NATS_URL            = "NATS_URL"
	NATS_HOST           = "NATS_HOST"
	NATS_PORT           = "NATS_PORT"
	NATS_USER           = "NATS_USER"
	NATS_PASSWORD       = "NATS_PASSWORD"
	NATS_STREAM         = "NATS_STREAM"
	NATS_SUBJECT_PREFIX = "NATS_SUBJECT_PREFIX"
	NATS_CONSUMER_NAME  = "NATS_CONSUMER_NAME"
	NATS_CONSUMER_BATCH = "NATS_CONSUMER_BATCH"
	NATS_CONSUMER_WAIT  = "NATS_CONSUMER_WAIT_SECONDS"
)

type NatsConf struct {
	url           string
	stream        string
	subjectPrefix string
	consumerName  string
	consumerBatch int
	consumerWait  time.Duration
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

	consumerName := os.Getenv(NATS_CONSUMER_NAME)
	if consumerName == "" {
		consumerName = "analytics_consumer"
	}

	consumerBatch := 100
	if rawBatch := os.Getenv(NATS_CONSUMER_BATCH); rawBatch != "" {
		batch, err := strconv.Atoi(rawBatch)
		if err != nil || batch <= 0 {
			return nil, errors.New("failed to get valid nats consumer batch")
		}
		consumerBatch = batch
	}

	consumerWait := 5 * time.Second
	if rawWait := os.Getenv(NATS_CONSUMER_WAIT); rawWait != "" {
		waitSeconds, err := strconv.Atoi(rawWait)
		if err != nil || waitSeconds <= 0 {
			return nil, errors.New("failed to get valid nats consumer wait seconds")
		}
		consumerWait = time.Duration(waitSeconds) * time.Second
	}

	return &NatsConf{
		url:           url,
		stream:        stream,
		subjectPrefix: subjectPrefix,
		consumerName:  consumerName,
		consumerBatch: consumerBatch,
		consumerWait:  consumerWait,
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

func (conf *NatsConf) ConsumerName() string {
	return conf.consumerName
}

func (conf *NatsConf) ConsumerBatch() int {
	return conf.consumerBatch
}

func (conf *NatsConf) ConsumerWait() time.Duration {
	return conf.consumerWait
}

package configs

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	KAFKA_BROKERS                   = "KAFKA_BROKERS"
	KAFKA_CLIENT_ID                 = "KAFKA_CLIENT_ID"
	KAFKA_TOPIC                     = "KAFKA_TOPIC"
	KAFKA_DLQ_TOPIC                 = "KAFKA_DLQ_TOPIC"
	KAFKA_CONSUMER_GROUP            = "KAFKA_CONSUMER_GROUP"
	KAFKA_CONSUMER_BATCH            = "KAFKA_CONSUMER_BATCH"
	KAFKA_CONSUMER_WAIT             = "KAFKA_CONSUMER_WAIT_SECONDS"
	KAFKA_HANDLER_MAX_ATTEMPTS      = "KAFKA_HANDLER_MAX_ATTEMPTS"
	KAFKA_RETRY_BACKOFF             = "KAFKA_RETRY_BACKOFF_SECONDS"
	KAFKA_RETRY_MAX_BACKOFF         = "KAFKA_RETRY_MAX_BACKOFF_SECONDS"
	KAFKA_PRODUCE_TIMEOUT           = "KAFKA_PRODUCE_TIMEOUT_SECONDS"
	KAFKA_COMMIT_TIMEOUT            = "KAFKA_COMMIT_TIMEOUT_SECONDS"
	defaultKafkaClientID            = "url_shortener"
	defaultKafkaConsumerGroup       = "analytics_consumer"
	defaultKafkaConsumerBatch       = 100
	defaultKafkaConsumerWaitSeconds = 5
	defaultKafkaHandlerMaxAttempts  = 5
	defaultKafkaRetryBackoffSeconds = 1
	defaultKafkaMaxBackoffSeconds   = 16
	defaultKafkaProduceTimeout      = 10
	defaultKafkaCommitTimeout       = 10
)

type KafkaConf struct {
	brokers            []string
	clientID           string
	topic              string
	dlqTopic           string
	consumerGroup      string
	consumerBatch      int
	consumerWait       time.Duration
	handlerMaxAttempts int
	retryBackoff       time.Duration
	retryMaxBackoff    time.Duration
	produceTimeout     time.Duration
	commitTimeout      time.Duration
}

func NewKafkaConf() (*KafkaConf, error) {
	brokers := splitAndTrim(os.Getenv(KAFKA_BROKERS))
	if len(brokers) == 0 {
		return nil, errors.New("failed to get kafka brokers")
	}

	topic := strings.TrimSpace(os.Getenv(KAFKA_TOPIC))
	if topic == "" {
		return nil, errors.New("failed to get kafka topic")
	}

	dlqTopic := strings.TrimSpace(os.Getenv(KAFKA_DLQ_TOPIC))
	if dlqTopic == "" {
		return nil, errors.New("failed to get kafka dlq topic")
	}

	clientID := strings.TrimSpace(os.Getenv(KAFKA_CLIENT_ID))
	if clientID == "" {
		clientID = defaultKafkaClientID
	}

	consumerGroup := strings.TrimSpace(os.Getenv(KAFKA_CONSUMER_GROUP))
	if consumerGroup == "" {
		consumerGroup = defaultKafkaConsumerGroup
	}

	consumerBatch, err := positiveEnvInt(KAFKA_CONSUMER_BATCH, defaultKafkaConsumerBatch)
	if err != nil {
		return nil, errors.New("failed to get valid kafka consumer batch")
	}

	consumerWaitSeconds, err := positiveEnvInt(KAFKA_CONSUMER_WAIT, defaultKafkaConsumerWaitSeconds)
	if err != nil {
		return nil, errors.New("failed to get valid kafka consumer wait seconds")
	}

	handlerMaxAttempts, err := positiveEnvInt(KAFKA_HANDLER_MAX_ATTEMPTS, defaultKafkaHandlerMaxAttempts)
	if err != nil {
		return nil, errors.New("failed to get valid kafka handler max attempts")
	}

	retryBackoffSeconds, err := positiveEnvInt(KAFKA_RETRY_BACKOFF, defaultKafkaRetryBackoffSeconds)
	if err != nil {
		return nil, errors.New("failed to get valid kafka retry backoff seconds")
	}

	retryMaxBackoffSeconds, err := positiveEnvInt(KAFKA_RETRY_MAX_BACKOFF, defaultKafkaMaxBackoffSeconds)
	if err != nil || retryMaxBackoffSeconds < retryBackoffSeconds {
		return nil, errors.New("failed to get valid kafka retry max backoff seconds")
	}

	produceTimeoutSeconds, err := positiveEnvInt(KAFKA_PRODUCE_TIMEOUT, defaultKafkaProduceTimeout)
	if err != nil {
		return nil, errors.New("failed to get valid kafka produce timeout seconds")
	}

	commitTimeoutSeconds, err := positiveEnvInt(KAFKA_COMMIT_TIMEOUT, defaultKafkaCommitTimeout)
	if err != nil {
		return nil, errors.New("failed to get valid kafka commit timeout seconds")
	}

	return &KafkaConf{
		brokers:            brokers,
		clientID:           clientID,
		topic:              topic,
		dlqTopic:           dlqTopic,
		consumerGroup:      consumerGroup,
		consumerBatch:      consumerBatch,
		consumerWait:       time.Duration(consumerWaitSeconds) * time.Second,
		handlerMaxAttempts: handlerMaxAttempts,
		retryBackoff:       time.Duration(retryBackoffSeconds) * time.Second,
		retryMaxBackoff:    time.Duration(retryMaxBackoffSeconds) * time.Second,
		produceTimeout:     time.Duration(produceTimeoutSeconds) * time.Second,
		commitTimeout:      time.Duration(commitTimeoutSeconds) * time.Second,
	}, nil
}

func splitAndTrim(raw string) []string {
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			values = append(values, value)
		}
	}
	return values
}

func positiveEnvInt(key string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, errors.New("value must be a positive integer")
	}
	return value, nil
}

func (conf *KafkaConf) Brokers() []string {
	return append([]string(nil), conf.brokers...)
}

func (conf *KafkaConf) ClientID() string {
	return conf.clientID
}

func (conf *KafkaConf) Topic() string {
	return conf.topic
}

func (conf *KafkaConf) DLQTopic() string {
	return conf.dlqTopic
}

func (conf *KafkaConf) ConsumerGroup() string {
	return conf.consumerGroup
}

func (conf *KafkaConf) ConsumerBatch() int {
	return conf.consumerBatch
}

func (conf *KafkaConf) ConsumerWait() time.Duration {
	return conf.consumerWait
}

func (conf *KafkaConf) HandlerMaxAttempts() int {
	return conf.handlerMaxAttempts
}

func (conf *KafkaConf) RetryBackoff() time.Duration {
	return conf.retryBackoff
}

func (conf *KafkaConf) RetryMaxBackoff() time.Duration {
	return conf.retryMaxBackoff
}

func (conf *KafkaConf) ProduceTimeout() time.Duration {
	return conf.produceTimeout
}

func (conf *KafkaConf) CommitTimeout() time.Duration {
	return conf.commitTimeout
}

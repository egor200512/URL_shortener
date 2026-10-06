package kafka

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
)

type Consumer struct {
	client          client
	dlqTopic        string
	batch           int
	maxAttempts     int
	retryBackoff    time.Duration
	retryMaxBackoff time.Duration
	commitTimeout   time.Duration
	now             func() time.Time
	closeOnce       sync.Once
}

type partitionID struct {
	topic     string
	partition int32
}

type partitionResult struct {
	lastHandled *record
	err         error
}

func NewConsumer(conf *configs.KafkaConf) (*Consumer, error) {
	kafkaClient, err := newConsumerClient(conf)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:          kafkaClient,
		dlqTopic:        conf.DLQTopic(),
		batch:           conf.ConsumerBatch(),
		maxAttempts:     conf.HandlerMaxAttempts(),
		retryBackoff:    conf.RetryBackoff(),
		retryMaxBackoff: conf.RetryMaxBackoff(),
		commitTimeout:   conf.CommitTimeout(),
		now:             time.Now,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, handler broker.MessageHandler) error {
	if handler == nil {
		return errors.New("kafka message handler is not initialized")
	}
	if c == nil || c.client == nil {
		return errors.New("kafka consumer is not initialized")
	}

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		records, err := c.client.Poll(ctx, c.batch)
		if err != nil {
			c.client.AllowRebalance()
			return err
		}

		err = c.processBatch(ctx, records, handler)
		c.client.AllowRebalance()
		if err != nil {
			return err
		}
	}
}

func (c *Consumer) processBatch(ctx context.Context, records []*record, handler broker.MessageHandler) error {
	if len(records) == 0 {
		return nil
	}

	partitions := make(map[partitionID][]*record)
	for _, message := range records {
		if message == nil {
			continue
		}
		id := partitionID{topic: message.topic, partition: message.partition}
		partitions[id] = append(partitions[id], message)
	}

	results := make(chan partitionResult, len(partitions))
	var wg sync.WaitGroup
	for _, messages := range partitions {
		sort.Slice(messages, func(i, j int) bool { return messages[i].offset < messages[j].offset })
		wg.Add(1)
		go func(partitionRecords []*record) {
			defer wg.Done()
			lastHandled, err := c.processPartition(ctx, partitionRecords, handler)
			results <- partitionResult{lastHandled: lastHandled, err: err}
		}(messages)
	}
	wg.Wait()
	close(results)

	commits := make([]*record, 0, len(partitions))
	errs := make([]error, 0)
	for result := range results {
		if result.lastHandled != nil {
			commits = append(commits, result.lastHandled)
		}
		if result.err != nil {
			errs = append(errs, result.err)
		}
	}

	if len(commits) > 0 {
		commitCtx, cancel := context.WithTimeout(context.Background(), c.commitTimeout)
		commitErr := c.client.Commit(commitCtx, commits...)
		cancel()
		if commitErr != nil {
			errs = append(errs, fmt.Errorf("commit kafka offsets: %w", commitErr))
		}
	}

	return errors.Join(errs...)
}

func (c *Consumer) processPartition(ctx context.Context, records []*record, handler broker.MessageHandler) (*record, error) {
	var lastHandled *record
	for _, message := range records {
		if err := c.handleWithRetry(ctx, message, handler); err != nil {
			return lastHandled, err
		}
		lastHandled = message
	}
	return lastHandled, nil
}

func (c *Consumer) handleWithRetry(ctx context.Context, message *record, handler broker.MessageHandler) error {
	var handlerErr error
	for attempt := 1; attempt <= c.maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		handlerErr = handler(ctx, message.value)
		if handlerErr == nil {
			return nil
		}

		slog.Warn("failed to handle kafka message",
			"topic", message.topic,
			"partition", message.partition,
			"offset", message.offset,
			"attempt", attempt,
			"max_attempts", c.maxAttempts,
			"error", handlerErr,
		)
		if attempt == c.maxAttempts {
			break
		}
		if err := waitForRetry(ctx, c.backoff(attempt)); err != nil {
			return err
		}
	}

	if err := c.publishToDLQ(ctx, message, handlerErr); err != nil {
		return fmt.Errorf("publish kafka message to dlq: %w", err)
	}
	slog.Error("kafka message moved to dlq",
		"topic", message.topic,
		"partition", message.partition,
		"offset", message.offset,
		"dlq_topic", c.dlqTopic,
		"error", handlerErr,
	)
	return nil
}

func (c *Consumer) publishToDLQ(ctx context.Context, source *record, handlerErr error) error {
	return c.client.Produce(ctx, &record{
		topic:     c.dlqTopic,
		partition: -1,
		key:       append([]byte(nil), source.key...),
		value:     append([]byte(nil), source.value...),
		timestamp: c.now().UTC(),
		headers: []header{
			{key: "source-topic", value: []byte(source.topic)},
			{key: "source-partition", value: []byte(strconv.FormatInt(int64(source.partition), 10))},
			{key: "source-offset", value: []byte(strconv.FormatInt(source.offset, 10))},
			{key: "handler-error", value: []byte(handlerErr.Error())},
			{key: "handler-attempts", value: []byte(strconv.Itoa(c.maxAttempts))},
			{key: "failed-at", value: []byte(c.now().UTC().Format(time.RFC3339Nano))},
		},
	})
}

func (c *Consumer) backoff(failedAttempt int) time.Duration {
	delay := c.retryBackoff
	for i := 1; i < failedAttempt && delay < c.retryMaxBackoff; i++ {
		if delay > c.retryMaxBackoff/2 {
			return c.retryMaxBackoff
		}
		delay *= 2
	}
	if delay > c.retryMaxBackoff {
		return c.retryMaxBackoff
	}
	return delay
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *Consumer) Close() {
	if c == nil || c.client == nil {
		return
	}
	c.closeOnce.Do(c.client.Close)
}

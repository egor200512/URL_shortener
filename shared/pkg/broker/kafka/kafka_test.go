package kafka

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/shared/pkg/broker"
	"github.com/stretchr/testify/require"
)

func TestProducerPublish(t *testing.T) {
	client := &fakeClient{}
	producer := &Producer{client: client, topic: "link_events"}

	err := producer.Publish(context.Background(), producer.Topic(), "abc123", []byte("payload"))

	require.NoError(t, err)
	require.Len(t, client.producedRecords(), 1)
	message := client.producedRecords()[0]
	require.Equal(t, "link_events", message.topic)
	require.Equal(t, []byte("abc123"), message.key)
	require.Equal(t, []byte("payload"), message.value)
}

func TestProducerValidationAndClose(t *testing.T) {
	require.ErrorContains(t, (*Producer)(nil).Publish(context.Background(), "topic", "key", nil), "not initialized")

	client := &fakeClient{}
	producer := &Producer{client: client, topic: "link_events"}
	require.ErrorContains(t, producer.Publish(context.Background(), "", "key", nil), "topic")
	require.ErrorContains(t, producer.Publish(context.Background(), "topic", "", nil), "key")

	producer.Close()
	producer.Close()
	require.Equal(t, 1, client.closeCount())
}

func TestConsumerCommitsAfterSuccessfulHandling(t *testing.T) {
	client := &fakeClient{
		pollResults: []pollResult{
			{records: []*record{
				testRecord(0, 10, "first"),
				testRecord(0, 11, "second"),
				testRecord(1, 3, "third"),
			}},
			{err: context.Canceled},
		},
	}
	consumer := testConsumer(client)

	var mu sync.Mutex
	handled := make(map[string]int)
	err := consumer.Consume(context.Background(), func(_ context.Context, payload []byte) error {
		mu.Lock()
		defer mu.Unlock()
		handled[string(payload)]++
		return nil
	})

	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, map[string]int{"first": 1, "second": 1, "third": 1}, handled)
	commits := client.committedRecords()
	require.Len(t, commits, 1)
	require.ElementsMatch(t, []int64{11, 3}, []int64{commits[0][0].offset, commits[0][1].offset})
	require.GreaterOrEqual(t, client.allowCount(), 2)
}

func TestConsumerProcessesPartitionsInParallelAndKeepsPartitionOrder(t *testing.T) {
	client := &fakeClient{pollResults: []pollResult{
		{records: []*record{
			testRecord(0, 11, "partition-zero-second"),
			testRecord(1, 3, "partition-one"),
			testRecord(0, 10, "partition-zero-first"),
		}},
		{err: context.Canceled},
	}}
	consumer := testConsumer(client)
	partitionZeroStarted := make(chan struct{})
	releasePartitionZero := make(chan struct{})
	partitionOneHandled := make(chan struct{})
	done := make(chan error, 1)

	var mu sync.Mutex
	var partitionZeroOrder []string
	go func() {
		done <- consumer.Consume(context.Background(), func(_ context.Context, payload []byte) error {
			value := string(payload)
			switch value {
			case "partition-zero-first":
				close(partitionZeroStarted)
				<-releasePartitionZero
			case "partition-one":
				close(partitionOneHandled)
			}
			if value != "partition-one" {
				mu.Lock()
				partitionZeroOrder = append(partitionZeroOrder, value)
				mu.Unlock()
			}
			return nil
		})
	}()

	select {
	case <-partitionZeroStarted:
	case <-time.After(time.Second):
		t.Fatal("partition zero did not start")
	}
	select {
	case <-partitionOneHandled:
	case <-time.After(time.Second):
		t.Fatal("partition one was blocked by partition zero")
	}
	close(releasePartitionZero)

	require.ErrorIs(t, <-done, context.Canceled)
	mu.Lock()
	require.Equal(t, []string{"partition-zero-first", "partition-zero-second"}, partitionZeroOrder)
	mu.Unlock()
}

func TestConsumerRetriesThenCommits(t *testing.T) {
	client := &fakeClient{pollResults: []pollResult{
		{records: []*record{testRecord(0, 10, "payload")}},
		{err: context.Canceled},
	}}
	consumer := testConsumer(client)
	consumer.maxAttempts = 3

	attempts := 0
	err := consumer.Consume(context.Background(), func(context.Context, []byte) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 3, attempts)
	require.Empty(t, client.producedRecords())
	require.Equal(t, int64(10), client.committedRecords()[0][0].offset)
}

func TestConsumerMovesPoisonMessageToDLQBeforeCommit(t *testing.T) {
	client := &fakeClient{pollResults: []pollResult{
		{records: []*record{testRecord(2, 42, "bad-payload")}},
		{err: context.Canceled},
	}}
	consumer := testConsumer(client)
	consumer.maxAttempts = 3

	attempts := 0
	err := consumer.Consume(context.Background(), func(context.Context, []byte) error {
		attempts++
		return errors.New("permanent error")
	})

	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 3, attempts)
	require.Len(t, client.producedRecords(), 1)
	dlqMessage := client.producedRecords()[0]
	require.Equal(t, "link_events.dlq", dlqMessage.topic)
	require.Equal(t, []byte("abc123"), dlqMessage.key)
	require.Equal(t, []byte("bad-payload"), dlqMessage.value)
	require.Equal(t, "link_events", headerValue(dlqMessage, "source-topic"))
	require.Equal(t, "2", headerValue(dlqMessage, "source-partition"))
	require.Equal(t, "42", headerValue(dlqMessage, "source-offset"))
	require.Equal(t, "3", headerValue(dlqMessage, "handler-attempts"))
	require.Equal(t, int64(42), client.committedRecords()[0][0].offset)
}

func TestConsumerDoesNotCommitFailedMessageWhenDLQPublishFails(t *testing.T) {
	client := &fakeClient{
		pollResults: []pollResult{{records: []*record{
			testRecord(0, 10, "ok"),
			testRecord(0, 11, "bad"),
			testRecord(0, 12, "must-not-run"),
		}}},
		produceErr: errors.New("kafka unavailable"),
	}
	consumer := testConsumer(client)
	consumer.maxAttempts = 1

	var handled []string
	err := consumer.Consume(context.Background(), func(_ context.Context, payload []byte) error {
		handled = append(handled, string(payload))
		if string(payload) == "bad" {
			return errors.New("handler failed")
		}
		return nil
	})

	require.ErrorContains(t, err, "publish kafka message to dlq")
	require.Equal(t, []string{"ok", "bad"}, handled)
	commits := client.committedRecords()
	require.Len(t, commits, 1)
	require.Len(t, commits[0], 1)
	require.Equal(t, int64(10), commits[0][0].offset)
}

func TestConsumerCommitErrorIsReturned(t *testing.T) {
	client := &fakeClient{
		pollResults: []pollResult{{records: []*record{testRecord(0, 10, "payload")}}},
		commitErr:   errors.New("commit failed"),
	}
	consumer := testConsumer(client)

	err := consumer.Consume(context.Background(), func(context.Context, []byte) error { return nil })

	require.ErrorContains(t, err, "commit kafka offsets")
}

func TestConsumerShutdownDoesNotHang(t *testing.T) {
	client := &fakeClient{blockPoll: true}
	consumer := testConsumer(client)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- consumer.Consume(ctx, func(context.Context, []byte) error { return nil })
	}()
	cancel()

	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("consumer did not stop after context cancellation")
	}
}

func TestConsumerValidationAndClose(t *testing.T) {
	require.ErrorContains(t, (*Consumer)(nil).Consume(context.Background(), func(context.Context, []byte) error { return nil }), "not initialized")
	require.ErrorContains(t, (&Consumer{}).Consume(context.Background(), nil), "handler")

	client := &fakeClient{}
	consumer := testConsumer(client)
	consumer.Close()
	consumer.Close()
	require.Equal(t, 1, client.closeCount())
}

func testConsumer(kafkaClient client) *Consumer {
	return &Consumer{
		client:          kafkaClient,
		dlqTopic:        "link_events.dlq",
		batch:           100,
		maxAttempts:     5,
		retryBackoff:    time.Millisecond,
		retryMaxBackoff: 4 * time.Millisecond,
		commitTimeout:   time.Second,
		now: func() time.Time {
			return time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
		},
	}
}

func testRecord(partition int32, offset int64, payload string) *record {
	return &record{
		topic:     "link_events",
		partition: partition,
		offset:    offset,
		key:       []byte("abc123"),
		value:     []byte(payload),
	}
}

func headerValue(message *record, key string) string {
	for _, messageHeader := range message.headers {
		if messageHeader.key == key {
			return string(messageHeader.value)
		}
	}
	return ""
}

type pollResult struct {
	records []*record
	err     error
}

type fakeClient struct {
	mu          sync.Mutex
	pollResults []pollResult
	pollIndex   int
	blockPoll   bool
	produceErr  error
	commitErr   error
	produced    []*record
	committed   [][]*record
	allows      int
	closes      int
}

func (c *fakeClient) Poll(ctx context.Context, _ int) ([]*record, error) {
	c.mu.Lock()
	if c.pollIndex < len(c.pollResults) {
		result := c.pollResults[c.pollIndex]
		c.pollIndex++
		c.mu.Unlock()
		return result.records, result.err
	}
	block := c.blockPoll
	c.mu.Unlock()

	if block {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return nil, context.Canceled
}

func (c *fakeClient) Produce(_ context.Context, message *record) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.produced = append(c.produced, message)
	return c.produceErr
}

func (c *fakeClient) Commit(_ context.Context, records ...*record) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.committed = append(c.committed, append([]*record(nil), records...))
	return c.commitErr
}

func (c *fakeClient) AllowRebalance() {
	c.mu.Lock()
	c.allows++
	c.mu.Unlock()
}

func (c *fakeClient) Close() {
	c.mu.Lock()
	c.closes++
	c.mu.Unlock()
}

func (c *fakeClient) producedRecords() []*record {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]*record(nil), c.produced...)
}

func (c *fakeClient) committedRecords() [][]*record {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make([][]*record, len(c.committed))
	for i := range c.committed {
		result[i] = append([]*record(nil), c.committed[i]...)
	}
	return result
}

func (c *fakeClient) allowCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.allows
}

func (c *fakeClient) closeCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closes
}

var _ broker.MessageHandler = func(context.Context, []byte) error { return nil }

package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/twmb/franz-go/pkg/kgo"
)

var errClientClosed = errors.New("kafka client is closed")

type header struct {
	key   string
	value []byte
}

type record struct {
	topic     string
	partition int32
	offset    int64
	key       []byte
	value     []byte
	timestamp time.Time
	headers   []header
	native    *kgo.Record
}

type client interface {
	Poll(context.Context, int) ([]*record, error)
	Produce(context.Context, *record) error
	Commit(context.Context, ...*record) error
	AllowRebalance()
	Close()
}

type franzClient struct {
	client *kgo.Client
}

func newProducerClient(conf *configs.KafkaConf) (client, error) {
	return newFranzClient(conf, false)
}

func newConsumerClient(conf *configs.KafkaConf) (client, error) {
	return newFranzClient(conf, true)
}

func newFranzClient(conf *configs.KafkaConf, consumer bool) (client, error) {
	if conf == nil {
		return nil, errors.New("kafka config is not initialized")
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(conf.Brokers()...),
		kgo.ClientID(conf.ClientID()),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordDeliveryTimeout(conf.ProduceTimeout()),
		kgo.RecordPartitioner(kgo.StickyKeyPartitioner(nil)),
	}
	if consumer {
		opts = append(opts,
			kgo.ConsumerGroup(conf.ConsumerGroup()),
			kgo.ConsumeTopics(conf.Topic()),
			kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
			kgo.DisableAutoCommit(),
			kgo.BlockRebalanceOnPoll(),
			kgo.FetchMaxWait(conf.ConsumerWait()),
		)
	}

	kafkaClient, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("create kafka client: %w", err)
	}
	return &franzClient{client: kafkaClient}, nil
}

func (c *franzClient) Poll(ctx context.Context, maxRecords int) ([]*record, error) {
	if c == nil || c.client == nil {
		return nil, errClientClosed
	}

	fetches := c.client.PollRecords(ctx, maxRecords)
	if fetches.IsClientClosed() {
		return nil, errClientClosed
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if fetchErrors := fetches.Errors(); len(fetchErrors) > 0 {
		errs := make([]error, 0, len(fetchErrors))
		for _, fetchErr := range fetchErrors {
			errs = append(errs, fmt.Errorf("fetch %s[%d]: %w", fetchErr.Topic, fetchErr.Partition, fetchErr.Err))
		}
		return nil, errors.Join(errs...)
	}

	kafkaRecords := fetches.Records()
	records := make([]*record, 0, len(kafkaRecords))
	for _, kafkaRecord := range kafkaRecords {
		records = append(records, recordFromKafka(kafkaRecord))
	}
	return records, nil
}

func (c *franzClient) Produce(ctx context.Context, message *record) error {
	if c == nil || c.client == nil {
		return errClientClosed
	}
	if message == nil {
		return errors.New("kafka record is not initialized")
	}

	return c.client.ProduceSync(ctx, recordToKafka(message)).FirstErr()
}

func (c *franzClient) Commit(ctx context.Context, records ...*record) error {
	if c == nil || c.client == nil {
		return errClientClosed
	}

	kafkaRecords := make([]*kgo.Record, 0, len(records))
	for _, message := range records {
		if message == nil {
			continue
		}
		kafkaRecords = append(kafkaRecords, recordToKafka(message))
	}
	if len(kafkaRecords) == 0 {
		return nil
	}
	return c.client.CommitRecords(ctx, kafkaRecords...)
}

func (c *franzClient) AllowRebalance() {
	if c != nil && c.client != nil {
		c.client.AllowRebalance()
	}
}

func (c *franzClient) Close() {
	if c != nil && c.client != nil {
		c.client.CloseAllowingRebalance()
	}
}

func recordFromKafka(source *kgo.Record) *record {
	message := &record{
		topic:     source.Topic,
		partition: source.Partition,
		offset:    source.Offset,
		key:       append([]byte(nil), source.Key...),
		value:     append([]byte(nil), source.Value...),
		timestamp: source.Timestamp,
		native:    source,
	}
	for _, sourceHeader := range source.Headers {
		message.headers = append(message.headers, header{
			key:   sourceHeader.Key,
			value: append([]byte(nil), sourceHeader.Value...),
		})
	}
	return message
}

func recordToKafka(message *record) *kgo.Record {
	if message.native != nil {
		return message.native
	}

	kafkaRecord := &kgo.Record{
		Topic:     message.topic,
		Key:       append([]byte(nil), message.key...),
		Value:     append([]byte(nil), message.value...),
		Timestamp: message.timestamp,
	}
	if message.partition >= 0 {
		kafkaRecord.Partition = message.partition
	}
	for _, messageHeader := range message.headers {
		kafkaRecord.Headers = append(kafkaRecord.Headers, kgo.RecordHeader{
			Key:   messageHeader.key,
			Value: append([]byte(nil), messageHeader.value...),
		})
	}
	return kafkaRecord
}

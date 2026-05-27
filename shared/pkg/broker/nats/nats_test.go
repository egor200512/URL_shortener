package nats

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	natsgo "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestProducer_Subject(t *testing.T) {
	producer := &Producer{subjectPrefix: "links.events"}

	require.Equal(t, "links.events.created", producer.Subject("created"))
}

func TestProducer_Close(t *testing.T) {
	require.NotPanics(t, func() {
		(*Producer)(nil).Close()
		(&Producer{}).Close()
	})
}

func TestProducer_Publish_NotInitialized(t *testing.T) {
	err := (*Producer)(nil).Publish(context.Background(), "links.events.created", []byte("payload"))

	require.Error(t, err)
	require.Contains(t, err.Error(), "not initialized")
}

func TestConsumer_Subject(t *testing.T) {
	consumer := &Consumer{subjectPrefix: "links.events"}

	require.Equal(t, "links.events.created", consumer.Subject("created"))
}

func TestConsumer_Close(t *testing.T) {
	require.NotPanics(t, func() {
		(*Consumer)(nil).Close()
		(&Consumer{}).Close()
	})
}

func TestConsumer_Consume_NotInitialized(t *testing.T) {
	err := (*Consumer)(nil).Consume(context.Background(), func(context.Context, []byte) error {
		return nil
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "not initialized")
}

func TestConsumer_Consume_NilHandler(t *testing.T) {
	err := (&Consumer{}).Consume(context.Background(), nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "handler")
}

func TestConsumer_Consume_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	consumer := &Consumer{
		subscription: &natsgo.Subscription{},
		batch:        1,
		wait:         time.Millisecond,
	}

	err := consumer.Consume(ctx, func(context.Context, []byte) error {
		return nil
	})

	require.ErrorIs(t, err, context.Canceled)
}

func TestConsumer_HandleMessages_ContinuesAfterHandlerError(t *testing.T) {
	consumer := &Consumer{}
	messages := []*natsgo.Msg{
		natsgo.NewMsg("links.events.created"),
		natsgo.NewMsg("links.events.deleted"),
	}
	messages[0].Data = []byte("bad")
	messages[1].Data = []byte("ok")

	var handled [][]byte
	consumer.handleMessages(context.Background(), messages, func(_ context.Context, payload []byte) error {
		handled = append(handled, payload)
		if string(payload) == "bad" {
			return errors.New("handler error")
		}
		return nil
	})

	require.Equal(t, [][]byte{[]byte("bad"), []byte("ok")}, handled)
}

func TestConsumerErrorsAreSpecific(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "nil handler", err: (&Consumer{}).Consume(context.Background(), nil), want: "handler"},
		{name: "nil consumer", err: (*Consumer)(nil).Consume(context.Background(), func(context.Context, []byte) error { return nil }), want: "consumer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, tt.err)
			require.True(t, strings.Contains(tt.err.Error(), tt.want), tt.err.Error())
		})
	}
}

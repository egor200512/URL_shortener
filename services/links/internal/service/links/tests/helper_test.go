package links

import "context"

type fakeProducer struct{}

func (fakeProducer) Publish(context.Context, string, []byte) error {
	return nil
}

func (fakeProducer) Close() {}

func (fakeProducer) Prefix() string {
	return "links"
}

package consumer

import (
	"fmt"
	"time"

	gonats "github.com/nats-io/nats.go"
)

func (c *consumer) Fetch(batch int, wait time.Duration) ([]*gonats.Msg, error) {
	if c.sub == nil {
		return nil, fmt.Errorf("subscription is nil")
	}
	return c.sub.Fetch(batch, gonats.MaxWait(wait))
}

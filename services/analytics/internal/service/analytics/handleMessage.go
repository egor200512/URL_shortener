package analytics

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/egor200512/URL_shortener/shared/models"

	gonats "github.com/nats-io/nats.go"
)

func (service *analyticsService) HandleMessage(ctx context.Context, msg *gonats.Msg) error {
	if msg == nil {
		return nil
	}

	parts := strings.Split(msg.Subject, ".")
	eventType := parts[len(parts)-1]

	var u models.JetStreamUnit

	if err := json.Unmarshal(msg.Data, &u); err != nil {
		return err
	}

	payload := &models.LinkEvent{
		EventType:    eventType,
		UserID:       u.UserID,
		ShortLink:    u.ShortLink,
		OriginalLink: u.OriginalLink,
		ExecutedAt:   u.ExecutedAt,
	}

	return service.repo.InsertEvent(ctx, payload)
}

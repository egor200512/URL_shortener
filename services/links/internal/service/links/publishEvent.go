package links

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/egor200512/URL_shortener/shared/pkg/events"
	"github.com/google/uuid"
)

func (service *linksService) publishLinkEvent(ctx context.Context, eventType string, link *models.Link) {
	if service.producer == nil || link == nil {
		return
	}

	event := &events.LinkEvent{
		EventID:      uuid.NewString(),
		EventType:    eventType,
		UserID:       link.UserID.String(),
		ShortLink:    link.ShortLink,
		OriginalLink: link.OriginalLink,
		ExecutedAt:   time.Now(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		slog.Warn("failed to marshal link event", "event_type", eventType, "short_link", link.ShortLink, "error", err)
		return
	}

	if err := service.producer.Publish(ctx, service.producer.Subject(eventType), payload); err != nil {
		slog.Warn("failed to publish link event", "event_type", eventType, "short_link", link.ShortLink, "error", err)
	}
}

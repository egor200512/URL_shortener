package links

import (
	"context"
	"encoding/json"
	"log"
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
		log.Printf("failed to marshal link event: %s\n", err.Error())
		return
	}

	if err := service.producer.Publish(ctx, service.producer.Subject(eventType), payload); err != nil {
		log.Printf("failed to publish link event: %s\n", err.Error())
	}
}

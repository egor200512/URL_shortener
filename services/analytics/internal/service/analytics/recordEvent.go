package analytics

import (
	"context"
	"time"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/google/uuid"
)

func (service *analyticsService) RecordEvent(ctx context.Context, req *desc.RecordEventRequest) error {
	eventID := req.EventId
	if eventID == "" {
		eventID = uuid.NewString()
	}

	executedAt := time.Now()
	if req.ExecutedAt != nil {
		executedAt = req.ExecutedAt.AsTime()
	}

	return service.repo.InsertEvent(ctx, &models.LinkEvent{
		ID:           eventID,
		EventType:    req.EventType,
		UserID:       req.UserId,
		ShortLink:    req.ShortLink,
		OriginalLink: req.OriginalLink,
		ExecutedAt:   executedAt,
	})
}

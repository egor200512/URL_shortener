package analytics

import (
	"context"
	"time"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

func (service *analyticsService) RecordEvent(ctx context.Context, req *desc.RecordEventRequest) error {
	executedAt := time.Now()
	if req.ExecutedAt != nil {
		executedAt = req.ExecutedAt.AsTime()
	}

	return service.repo.InsertEvent(ctx, &models.LinkEvent{
		EventType:    req.EventType,
		UserID:       req.UserId,
		ShortLink:    req.ShortLink,
		OriginalLink: req.OriginalLink,
		ExecutedAt:   executedAt,
	})
}

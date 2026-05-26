package repository

import (
	"context"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
)

type IAnalyticsRepo interface {
	InsertEvent(ctx context.Context, event *models.LinkEvent) error
	GetEvents(ctx context.Context, limit, offset int32) ([]*models.LinkEvent, int32, error)
	GetLinkStats(ctx context.Context, shortLink string) (*models.LinkStats, error)
}

package repository

import (
	"context"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
)

type IAnalyticsRepo interface {
	GetLinkStats(ctx context.Context, shortLink string) (*models.LinkStats, error)
}

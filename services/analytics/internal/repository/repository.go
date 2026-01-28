package repository

import (
	"context"

	"github.com/egor200512/URL_shortener/shared/models"
)

type IAnalyticsRepo interface {
	InsertEvent(ctx context.Context, payload *models.LinkEvent) error
	GetEvents(ctx context.Context, limit, offset int32) ([]*models.LinkEvent, error)
}

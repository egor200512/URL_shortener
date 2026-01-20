package redis

import (
	"context"

	"github.com/egor200512/URL_shortener/services/links/models"
)

type ICache interface {
	SetShort(ctx context.Context, shortLink string, link *models.Link) error
	GetShort(ctx context.Context, shortLink string) (*models.Link, error)
	DelShort(ctx context.Context, shortLink string) error
}

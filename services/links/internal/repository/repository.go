package repository

import (
	"context"

	"github.com/egor200512/URL_shortener/services/links/models"
)

type ILinksRepo interface {
	GetByShortLink(ctx context.Context, shortLink string) (*models.Link, error)
	GetByOriginalLink(ctx context.Context, originalLink string) (*models.Link, error)
	InsertLink(ctx context.Context, req *models.CreateLinkReq) error
	DeleteLink(ctx context.Context, shortLink string, userID string) error
}

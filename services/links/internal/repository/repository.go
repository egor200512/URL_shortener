package repository

import (
	"context"

	"github.com/egor200512/URL_shortener/services/links/models"
)

type ILinksRepo interface {
	GetByShortLink(ctx context.Context, shortLink string) (*models.Link, error)
	GetByOriginalLink(ctx context.Context, originalLink string) (*models.Link, error)
	GetUserLinks(ctx context.Context, userID string, limit, offset int32) ([]string, int32, error)
	InsertLink(ctx context.Context, req *models.CreateLinkReq) (*models.Link, error)
	DeleteLink(ctx context.Context, shortLink string, userID string) error
}

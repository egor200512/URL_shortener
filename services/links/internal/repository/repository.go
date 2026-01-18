package repository

import (
	"context"
	"net/url"

	"github.com/egor200512/URL_shortener/services/links/models"
)

type ILinksRepo interface {
	CheckShortLink(ctx context.Context, shortLink string) (*models.Link, error)
	InsertLink(ctx context.Context, u *url.URL, shortLink string) error
}

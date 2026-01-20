package service

import (
	"context"
	"net/url"

	"github.com/egor200512/URL_shortener/services/links/models"
)

type ILinksService interface {
	CreateLink(context.Context, *url.URL) (string, error)
	GetLinkInfo(ctx context.Context, shortLink string) (*models.Link, error)
	GetUserLinks(ctx context.Context, limit, offset int32) ([]string, int32, error)
	DeleteLink(ctx context.Context, shortLink string) error
}

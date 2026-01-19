package service

import (
	"context"
	"net/url"

	"github.com/egor200512/URL_shortener/services/links/models"
)

type ILinksService interface {
	CreateLink(context.Context, *url.URL) (string, error)
	GetLinkInfo(ctx context.Context, shortLink string) (*models.Link, error)
}

package links

import (
	"context"

	"github.com/egor200512/URL_shortener/services/links/models"
)

func (service *linksService) GetLinkInfo(ctx context.Context, shortLink string) (*models.Link, error) {
	l, err := service.linksRepo.GetByShortLink(ctx, shortLink)
	if err != nil {
		return nil, err
	}

	if l == nil {
		return nil, nil
	}

	return l, nil
}

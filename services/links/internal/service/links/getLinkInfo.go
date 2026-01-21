package links

import (
	"context"
	"log"

	"github.com/egor200512/URL_shortener/services/links/models"
)

func (service *linksService) GetLinkInfo(ctx context.Context, shortLink string) (*models.Link, error) {
	l, err := service.cache.GetShort(ctx, shortLink)
	if err != nil {
		log.Printf("failed to get link from cache: %s\n", err.Error())
	}

	if l != nil {
		log.Println("link from cache")
		return l, nil
	}

	l, err = service.linksRepo.GetByShortLink(ctx, shortLink)
	if err != nil {
		return nil, err
	}

	if l == nil {
		return nil, nil
	}

	if err = service.cache.SetShort(ctx, shortLink, l); err != nil {
		log.Printf("failed to cache link: %s\n", err.Error())
	}

	return l, nil
}

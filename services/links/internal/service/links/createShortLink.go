package links

import (
	"context"
	"fmt"
	"net/url"

	l "github.com/egor200512/URL_shortener/shared/pkg/links"
)

func (service *linksService) CreateLink(ctx context.Context, u *url.URL) (string, error) {
	var shortLink string

	for {
		shortLink, err := l.GenerateShortLink()
		if err != nil {
			return "", fmt.Errorf("failed to generate short link: %s", err.Error())
		}

		l, err := service.linksRepo.CheckShortLink(ctx, shortLink)
		if err != nil {
			return "", err
		}

		if l == nil {
			break
		}
	}

	if err := service.linksRepo.InsertLink(ctx, u, shortLink); err != nil {
		return "", err
	}

	return shortLink, nil
}

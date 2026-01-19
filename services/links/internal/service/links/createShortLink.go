package links

import (
	"context"
	"fmt"
	"net/url"

	pkg "github.com/egor200512/URL_shortener/shared/pkg/links"
	"github.com/google/uuid"
)

func (service *linksService) CreateLink(ctx context.Context, u *url.URL) (string, error) {
	l, err := service.linksRepo.GetByOriginalLink(ctx, u.Host+u.Path)
	if err != nil {
		return "", err
	}

	if l != nil {
		return "", fmt.Errorf("%s alredy exists", u.Host+u.Path)
	}

	var shortLink string

	for {
		shortLink, err = pkg.GenerateShortLink()
		if err != nil {
			return "", err
		}

		l, err := service.linksRepo.GetByShortLink(ctx, shortLink)
		if err != nil {
			return "", err
		}

		if l == nil {
			break
		}
	}

	if err = service.linksRepo.InsertLink(ctx, u, shortLink, uuid.NewString()); err != nil {
		return "", err
	}

	return "", nil
}

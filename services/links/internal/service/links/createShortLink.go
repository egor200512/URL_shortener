package links

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/links"
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

	userID, ok := ctx.Value(jwt.ClaimsCtxKey).(string)
	if !ok {
		return "", errors.New("failed to get userID from ctx")
	}

	req := &models.CreateLinkReq{
		UserID:           userID,
		ShortLink:        shortLink,
		OriginalLinkHost: u.Host,
		OriginalLink:     u.Host + u.Path,
	}

	created, err := service.linksRepo.InsertLink(ctx, req)
	if err != nil {
		return "", err
	}

	_ = service.cache.SetShort(ctx, shortLink, created)

	return shortLink, nil
}

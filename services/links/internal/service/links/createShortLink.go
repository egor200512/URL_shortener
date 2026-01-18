package links

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	pkg "github.com/egor200512/URL_shortener/shared/pkg/links"

	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
)

func (service *linksService) CreateLink(ctx context.Context, u *url.URL) (string, error) {
	var shortLink string

	for {
		shortLink, err := pkg.GenerateShortLink()
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

	user_id, ok := ctx.Value(jwt.ClaimsCtxKey).(string)
	if !ok {
		return "", errors.New("failed to get user_id from ctx")
	}

	if err := service.linksRepo.InsertLink(ctx, u, shortLink, user_id); err != nil {
		return "", err
	}

	return shortLink, nil
}

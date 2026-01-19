package links

import (
	"context"
	"errors"
	"fmt"

	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
)

func (service *linksService) DeleteLink(ctx context.Context, shortLink string) error {
	l, err := service.linksRepo.GetByShortLink(ctx, shortLink)
	if err != nil {
		return err
	}

	if l == nil {
		return fmt.Errorf("link for %s doesn't exist", shortLink)
	}

	userID, ok := ctx.Value(jwt.ClaimsCtxKey).(string)
	if !ok {
		return errors.New("failed to get userID from ctx")
	}

	if err := service.linksRepo.DeleteLink(ctx, shortLink, userID); err != nil {
		return err
	}

	return nil
}

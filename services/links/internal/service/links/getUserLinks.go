package links

import (
	"context"
	"errors"
	"log"

	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
)

func (service *linksService) GetUserLinks(ctx context.Context, limit, offset int32) ([]string, int32, error) {
	userID, ok := ctx.Value(jwt.ClaimsCtxKey).(string)
	if !ok {
		log.Println("failed to get userID from ctx")
		return nil, 0, errors.New("failed to get userID from ctx")
	}

	return service.linksRepo.GetUserLinks(ctx, userID, limit, offset)
}

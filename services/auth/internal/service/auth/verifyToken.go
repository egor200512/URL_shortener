package auth

import (
	"context"

	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"github.com/google/uuid"
)

func (service *authService) VerifyToken(ctx context.Context, token string) (uuid.UUID, error) {
	claims, err := jwt.VerifyToken(token, []byte(service.jwtConf.Secret()))
	if err != nil {
		return uuid.Nil, err
	}

	return claims.ID, nil
}

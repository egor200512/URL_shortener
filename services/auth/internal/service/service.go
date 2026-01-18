package service

import (
	"context"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/google/uuid"
)

type IAuthService interface {
	Register(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string) (*models.AccessToken, error)
	VerifyToken(ctx context.Context, token string) (uuid.UUID, error)
}

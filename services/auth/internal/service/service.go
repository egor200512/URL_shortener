package service

import (
	"context"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
)

type IAuthService interface {
	Register(ctx context.Context, email string, password string) error
	Login(context.Context, string, string) (*models.AccessToken, error)
}

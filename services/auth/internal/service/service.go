package service

import (
	"context"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
)

type AuthService interface {
	Register(context.Context, *models.RegisterRequest) error
	Login(context.Context, string, string) (string, error)
}

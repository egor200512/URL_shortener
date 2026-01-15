package repository

import (
	"context"

	"github.com/egor200512/PassVault/internal/models"
)

type IAuthRepo interface {
	CheckRegistration(ctx context.Context, email string) (*models.User, error)
	InsertUser(ctx context.Context, regRequest *models.RegisterRequest) error
}

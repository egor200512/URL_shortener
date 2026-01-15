package service

import (
	"context"
)

type IAuthService interface {
	Register(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string) (string, error)
}

package auth

import (
	"context"

	r "github.com/egor200512/URL_shortener/services/auth/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	ID             = "id"
	USERS          = "users"
	EMAIL          = "email"
	SALT           = "salt"
	SALT_PASS_HASH = "salt_password_hash"
)

type authRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(ctx context.Context, pool *pgxpool.Pool) r.IAuthRepo {
	return &authRepo{
		pool: pool,
	}
}

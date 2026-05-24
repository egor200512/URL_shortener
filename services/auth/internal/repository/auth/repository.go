package auth

import (
	"context"

	r "github.com/egor200512/URL_shortener/services/auth/internal/repository"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	ID             = "id"
	USERS          = "auth.users"
	EMAIL          = "email"
	SALT           = "salt"
	SALT_PASS_HASH = "salt_password_hash"
)

type authRepo struct {
	pool dbPool
}

type dbPool interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func NewAuthRepo(ctx context.Context, pool *pgxpool.Pool) r.IAuthRepo {
	return &authRepo{
		pool: pool,
	}
}

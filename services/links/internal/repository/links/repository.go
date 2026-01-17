package links

import (
	"context"

	r "github.com/egor200512/URL_shortener/services/links/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type linksRepo struct {
	pool *pgxpool.Pool
}

func NewLinksRepo(ctx context.Context, pool *pgxpool.Pool) r.ILinksRepo {
	return &linksRepo{
		pool: pool,
	}
}

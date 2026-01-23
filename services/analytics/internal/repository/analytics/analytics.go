package analytics

import (
	"context"

	"github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type analyticsRepository struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepo(ctx context.Context, pool *pgxpool.Pool) repository.IAnalyticsRepo {
	return &analyticsRepository{
		pool: pool,
	}
}

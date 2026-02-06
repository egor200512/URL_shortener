package analytics

import (
	"context"

	"github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	EVENTS        = "analytics.link_events"
	ID            = "id"
	EVENT_TYPE    = "event_type"
	USER_ID       = "user_id"
	SHORT_LINK    = "short_link"
	ORIGINAL_LINK = "original_link"
	EXECUTED_AT   = "executed_at"
)

type analyticsRepository struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepo(ctx context.Context, pool *pgxpool.Pool) repository.IAnalyticsRepo {
	return &analyticsRepository{
		pool: pool,
	}
}

package analytics

import (
	"context"

	r "github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	LINK_EVENTS = "analytics.link_events"
	EVENT_TYPE  = "event_type"
	SHORT_LINK  = "short_link"
)

type analyticsRepository struct {
	pool dbPool
}

type dbPool interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func NewAnalyticsRepo(ctx context.Context, pool *pgxpool.Pool) r.IAnalyticsRepo {
	return &analyticsRepository{
		pool: pool,
	}
}

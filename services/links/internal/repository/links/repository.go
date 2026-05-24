package links

import (
	"context"

	r "github.com/egor200512/URL_shortener/services/links/internal/repository"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	LINKS              = "links.short_links"
	ID                 = "id"
	USER_ID            = "user_id"
	SHORT_LINK         = "short_link"
	ORIGINAL_LINK_HOST = "original_link_host"
	ORIGINAL_LINK      = "original_link"
	CREATED_AT         = "created_at"
)

type linksRepository struct {
	pool dbPool
}

type dbPool interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func NewLinksRepo(ctx context.Context, pool *pgxpool.Pool) r.ILinksRepo {
	return &linksRepository{
		pool: pool,
	}
}

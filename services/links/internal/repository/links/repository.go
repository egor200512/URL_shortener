package links

import (
	"context"

	r "github.com/egor200512/URL_shortener/services/links/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	LINKS              = "links.short_links"
	ID                 = "id"
	USER_ID            = "user_id"
	SHORT_LINK         = "short_code"
	ORIGINAL_LINK_HOST = "original_url_host"
	ORIGINAL_LINK      = "original_url"
	CREATED_AT         = "created_at"
)

type linksRepo struct {
	pool *pgxpool.Pool
}

func NewLinksRepo(ctx context.Context, pool *pgxpool.Pool) r.ILinksRepo {
	return &linksRepo{
		pool: pool,
	}
}

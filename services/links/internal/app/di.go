package app

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"

	lh "github.com/egor200512/URL_shortener/services/links/internal/handler"
	lr "github.com/egor200512/URL_shortener/services/links/internal/repository"
	lri "github.com/egor200512/URL_shortener/services/links/internal/repository/links"
	ls "github.com/egor200512/URL_shortener/services/links/internal/service"
	lsi "github.com/egor200512/URL_shortener/services/links/internal/service/links"
	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ProvideHttpConf() (*configs.HttpConf, error) {
	return configs.NewHttpConf()
}

func ProvideGrpcConf() (*configs.GrpcConf, error) {
	return configs.NewGRPCConf()
}

func ProvidePgConf() (*configs.PgConf, error) {
	return configs.NewPgConf()
}

func ProvideJwtConf() (configs.IJwtConf, error) {
	return configs.NewJwtConf()
}

func ProvidePgPool(ctx context.Context, cfg *configs.PgConf) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, cfg.LinksDSN())
}

func ProvideLinksRepo(ctx context.Context, pool *pgxpool.Pool) lr.ILinksRepo {
	return lri.NewLinksRepo(ctx, pool)
}

func ProvideLinksService(repo lr.ILinksRepo, jwtCfg configs.IJwtConf) ls.ILinksService {
	return lsi.NewLinksService(repo, jwtCfg)
}

func ProvideLinksHandler(svc ls.ILinksService) desc.LinksServiceServer {
	return lh.NewLinksRouter(svc)
}

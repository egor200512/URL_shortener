package app

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"

	ah "github.com/egor200512/URL_shortener/services/analytics/internal/handler"
	ar "github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	ari "github.com/egor200512/URL_shortener/services/analytics/internal/repository/analytics"
	as "github.com/egor200512/URL_shortener/services/analytics/internal/service"
	asi "github.com/egor200512/URL_shortener/services/analytics/internal/service/analytics"
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

func ProvidePgPool(ctx context.Context, cfg *configs.PgConf) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, cfg.AnalyticsDSN())
}

func ProvideAnalyticsRepo(ctx context.Context, pool *pgxpool.Pool) ar.IAnalyticsRepo {
	return ari.NewAnalyticsRepo(ctx, pool)
}

func ProvideAnalyticsService(repo ar.IAnalyticsRepo) as.IAnalyticsService {
	return asi.NewAnalyticsService(repo)
}

func ProvideAnalyticsHandler(svc as.IAnalyticsService) desc.AnalyticsServiceServer {
	return ah.NewAnalyticsRouter(svc)
}

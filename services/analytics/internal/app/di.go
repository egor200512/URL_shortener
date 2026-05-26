package app

import (
	"context"

	ah "github.com/egor200512/URL_shortener/services/analytics/internal/handler"
	ar "github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	ari "github.com/egor200512/URL_shortener/services/analytics/internal/repository/analytics"
	as "github.com/egor200512/URL_shortener/services/analytics/internal/service"
	asi "github.com/egor200512/URL_shortener/services/analytics/internal/service/analytics"
	"github.com/egor200512/URL_shortener/shared/configs"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
	nats "github.com/egor200512/URL_shortener/shared/pkg/broker/nats"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ProvideHttpConf() (*configs.HttpConf, error) {
	return configs.NewHttpConf()
}

func ProvideGrpcConf() (*configs.GrpcConf, error) {
	return configs.NewGRPCConf()
}

func ProvideMetricsConf() (*configs.MetricsConf, error) {
	return configs.NewMetricsConf()
}

func ProvidePgConf() (*configs.PgConf, error) {
	return configs.NewPgConf()
}

func ProvideNatsConf() (*configs.NatsConf, error) {
	return configs.NewNatsConf()
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

func ProvideConsumer(cfg *configs.NatsConf) broker.IConsumer {
	return nats.NewConsumer(cfg)
}

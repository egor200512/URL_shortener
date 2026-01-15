package app

import (
	"context"

	desc "github.com/egor200512/URL_shortener/services/auth/internal/gen_auth"

	ah "github.com/egor200512/URL_shortener/services/auth/internal/handler"
	ar "github.com/egor200512/URL_shortener/services/auth/internal/repository"
	ari "github.com/egor200512/URL_shortener/services/auth/internal/repository/auth"
	as "github.com/egor200512/URL_shortener/services/auth/internal/service"
	asi "github.com/egor200512/URL_shortener/services/auth/internal/service/auth"
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
	return pgxpool.Connect(ctx, cfg.AuthDSN())
}

func ProvideAuthRepo(ctx context.Context, pool *pgxpool.Pool) ar.IAuthRepo {
	return ari.NewAuthRepo(ctx, pool)
}

func ProvideAuthService(repo ar.IAuthRepo, jwtCfg configs.IJwtConf) as.IAuthService {
	return asi.NewAuthService(repo, jwtCfg)
}

func ProvideAuthHandler(svc as.IAuthService) desc.AuthServiceServer {
	return ah.NewAuthRouter(svc)
}

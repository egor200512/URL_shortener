package app

import (
	"context"
	"log"

	descA "github.com/egor200512/URL_shortener/shared/gen/auth"
	descL "github.com/egor200512/URL_shortener/shared/gen/links"
	cache "github.com/egor200512/URL_shortener/shared/pkg/cache"
	red "github.com/egor200512/URL_shortener/shared/pkg/cache/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

func ProvideRedisConf() (*configs.RedisConf, error) {
	return configs.NewRedisConf()
}

func ProvidePgPool(ctx context.Context, cfg *configs.PgConf) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, cfg.LinksDSN())
}

func ProvideLinksRepo(ctx context.Context, pool *pgxpool.Pool) lr.ILinksRepo {
	return lri.NewLinksRepo(ctx, pool)
}

func ProvideLinksService(repo lr.ILinksRepo, cache cache.ICache, cfg configs.IJwtConf) ls.ILinksService {
	return lsi.NewLinksService(repo, cache, cfg)
}

func ProvideLinksHandler(svc ls.ILinksService) descL.LinksServiceServer {
	return lh.NewLinksRouter(svc)
}

func ProvideCacheCli(cfg *configs.RedisConf) cache.ICache {
	return red.NewRedisCli(cfg)
}

func ProvideAuthServiceClient(conf *configs.GrpcConf) descA.AuthServiceClient {
	conn, err := grpc.NewClient(conf.AuthAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create auth client: %s", err.Error())
	}
	return descA.NewAuthServiceClient(conn)
}

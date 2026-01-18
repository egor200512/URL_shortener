package app

import (
	"context"
	"log"

	descA "github.com/egor200512/URL_shortener/shared/gen/auth"
	descL "github.com/egor200512/URL_shortener/shared/gen/links"
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

func ProvidePgPool(ctx context.Context, cfg *configs.PgConf) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, cfg.LinksDSN())
}

func ProvideLinksRepo(ctx context.Context, pool *pgxpool.Pool) lr.ILinksRepo {
	return lri.NewLinksRepo(ctx, pool)
}

func ProvideLinksService(repo lr.ILinksRepo, jwtCfg configs.IJwtConf) ls.ILinksService {
	return lsi.NewLinksService(repo, jwtCfg)
}

func ProvideLinksHandler(svc ls.ILinksService) descL.LinksServiceServer {
	return lh.NewLinksRouter(svc)
}

func ProvideAuthServiceClient(conf *configs.GrpcConf) descA.AuthServiceClient {
	conn, err := grpc.NewClient(conf.AuthAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create auth client: %s", err.Error())
	}
	return descA.NewAuthServiceClient(conn)
}

//go:build wireinject
// +build wireinject

package app

import (
	"context"

	"github.com/google/wire"
)

func InitServerApp(ctx context.Context) (*App, func(), error) {

	wire.Build(
		ProvideHttpConf,
		ProvideGrpcConf,
		ProvidePgConf,
		ProvideJwtConf,
		ProvideRedisConf,
		ProvideKafkaConf,

		ProvidePgPool,

		ProvideLinksRepo,

		ProvideLinksService,

		ProvideLinksHandler,

		ProvideAuthServiceClient,

		ProvideCacheCli,
		ProvideProducer,

		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}

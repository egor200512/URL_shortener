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

		ProvidePgPool,

		ProvideLinksRepo,

		ProvideLinksService,

		ProvideLinksHandler,

		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}

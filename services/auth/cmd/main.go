package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/egor200512/URL_shortener/services/auth/internal/app"
	"github.com/egor200512/URL_shortener/shared/configs"
	"github.com/egor200512/URL_shortener/shared/pkg/logger"
)

const ENV_FILE_PATH = "/home/egoor/url_shortener/.env"

func main() {
	logger.Configure("auth")
	configs.LoadConfig(ENV_FILE_PATH)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, cleanup, err := app.InitServerApp(ctx)
	if err != nil {
		logger.Fatal("failed to init auth app", err)
	}
	defer cleanup()

	if err := a.Run(ctx); err != nil {
		logger.Fatal("failed to run auth app", err)
	}
}

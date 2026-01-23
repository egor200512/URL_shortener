package main

import (
	"context"
	"log"

	"github.com/egor200512/URL_shortener/services/auth/internal/app"
	"github.com/egor200512/URL_shortener/shared/configs"
)

const ENV_FILE_PATH = "/home/egoor/url_shortener/.env"

func main() {
	configs.LoadConfig(ENV_FILE_PATH)

	ctx := context.Background()

	a, cleanup, err := app.InitServerApp(ctx)
	if err != nil {
		log.Fatalf("failed to init auth app: %v", err)
	}
	defer cleanup()

	if err := a.Run(ctx); err != nil {
		log.Fatalf("failed to run auth app: %v", err)
	}
}

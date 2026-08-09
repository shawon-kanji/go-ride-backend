package main

import (
	"context"
	"log"

	"go-ride-backend/internal/bootstrap"
	"go-ride-backend/internal/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app, err := bootstrap.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}

	if err := app.Router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

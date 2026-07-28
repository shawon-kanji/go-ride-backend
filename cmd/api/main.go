package main

import (
	"context"
	"log"

	"go-ride-backend/internal/bootstrap"
	"go-ride-backend/internal/config"
)

func main() {
	cfg, err := config.Load(context.Background())
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app, err := bootstrap.Build(cfg)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}

	if err := app.Router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

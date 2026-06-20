package bootstrap

import (
	"fmt"

	appuser "go-ride-backend/application/user"
	"go-ride-backend/infrastructure/db"
	"go-ride-backend/infrastructure/db/models"
	"go-ride-backend/infrastructure/repository"
	"go-ride-backend/infrastructure/security"
	"go-ride-backend/interfaces/http/handlers"
	"go-ride-backend/interfaces/http/routes"
	"go-ride-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
}

func Build(cfg *config.Config) (*App, error) {
	gormDB, err := db.NewGorm(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	if err := gormDB.AutoMigrate(&models.UserModel{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	userRepo := repository.NewUserRepositoryGorm(gormDB)
	hasher := security.NewBcryptHasher()
	jwtManager := security.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiryMinutes)

	signupUseCase := appuser.NewSignupUseCase(userRepo, hasher)
	loginUseCase := appuser.NewLoginUseCase(userRepo, hasher, jwtManager)
	authHandler := handlers.NewAuthHandler(signupUseCase, loginUseCase)

	router := routes.NewRouter(authHandler, jwtManager)
	return &App{Router: router}, nil
}

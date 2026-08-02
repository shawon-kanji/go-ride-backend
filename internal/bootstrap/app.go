package bootstrap

import (
	"fmt"

	appdriver "go-ride-backend/application/driver"
	appuser "go-ride-backend/application/user"
	appvehicle "go-ride-backend/application/vehicle"
	"go-ride-backend/infrastructure/db"
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

	userRepo := repository.NewUserRepositoryGorm(gormDB)
	driverRepo := repository.NewDriverRepositoryGorm(gormDB)
	vehicleRepo := repository.NewVehicleRepositoryGorm(gormDB)
	hasher := security.NewBcryptHasher()
	jwtManager := security.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiryMinutes, cfg.JWT.Issuer, cfg.JWT.Audience, cfg.JWT.DriverAudience)

	signupUseCase := appuser.NewSignupUseCase(userRepo, hasher)
	loginUseCase := appuser.NewLoginUseCase(userRepo, hasher, jwtManager)
	getProfileUseCase := appuser.NewGetProfileUseCase(userRepo)
	updateProfileUseCase := appuser.NewUpdateProfileUseCase(userRepo)
	changePasswordUseCase := appuser.NewChangePasswordUseCase(userRepo, hasher)
	deactivateAccountUseCase := appuser.NewDeactivateAccountUseCase(userRepo)
	authHandler := handlers.NewAuthHandler(
		signupUseCase,
		loginUseCase,
		getProfileUseCase,
		updateProfileUseCase,
		changePasswordUseCase,
		deactivateAccountUseCase,
	)
	driverSignupUseCase := appdriver.NewSignupUseCase(driverRepo, hasher)
	driverLoginUseCase := appdriver.NewLoginUseCase(driverRepo, hasher, jwtManager)
	driverGetProfileUseCase := appdriver.NewGetProfileUseCase(driverRepo)
	driverUpdateProfileUseCase := appdriver.NewUpdateProfileUseCase(driverRepo)
	driverUpdateOnlineStatusUseCase := appdriver.NewUpdateOnlineStatusUseCase(driverRepo, vehicleRepo)
	driverHandler := handlers.NewDriverHandler(
		driverSignupUseCase,
		driverLoginUseCase,
		driverGetProfileUseCase,
		driverUpdateProfileUseCase,
		driverUpdateOnlineStatusUseCase,
	)

	vehicleRegisterUseCase := appvehicle.NewRegisterUseCase(vehicleRepo)
	vehicleListUseCase := appvehicle.NewListUseCase(vehicleRepo)
	vehicleGetUseCase := appvehicle.NewGetUseCase(vehicleRepo)
	vehicleUpdateUseCase := appvehicle.NewUpdateUseCase(vehicleRepo)
	vehicleActivateUseCase := appvehicle.NewActivateUseCase(vehicleRepo)
	vehicleDeleteUseCase := appvehicle.NewDeleteUseCase(vehicleRepo)
	vehicleHandler := handlers.NewVehicleHandler(
		vehicleRegisterUseCase,
		vehicleListUseCase,
		vehicleGetUseCase,
		vehicleUpdateUseCase,
		vehicleActivateUseCase,
		vehicleDeleteUseCase,
	)

	router := routes.NewRouter(authHandler, driverHandler, vehicleHandler, jwtManager)
	return &App{Router: router}, nil
}

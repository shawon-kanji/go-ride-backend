package bootstrap

import (
	"context"
	"fmt"

	appdriver "go-ride-backend/application/driver"
	appkyc "go-ride-backend/application/kyc"
	appplaces "go-ride-backend/application/places"
	appsavedplace "go-ride-backend/application/savedplace"
	appuser "go-ride-backend/application/user"
	appvehicle "go-ride-backend/application/vehicle"
	"go-ride-backend/infrastructure/db"
	"go-ride-backend/infrastructure/maps"
	"go-ride-backend/infrastructure/repository"
	"go-ride-backend/infrastructure/security"
	"go-ride-backend/infrastructure/storage"
	"go-ride-backend/interfaces/http/handlers"
	"go-ride-backend/interfaces/http/routes"
	"go-ride-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
}

func Build(ctx context.Context, cfg *config.Config) (*App, error) {
	gormDB, err := db.NewGorm(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	userRepo := repository.NewUserRepositoryGorm(gormDB)
	driverRepo := repository.NewDriverRepositoryGorm(gormDB)
	vehicleRepo := repository.NewVehicleRepositoryGorm(gormDB)
	kycRepo := repository.NewKycRepositoryGorm(gormDB)
	savedPlaceRepo := repository.NewSavedPlaceRepositoryGorm(gormDB)
	hasher := security.NewBcryptHasher()
	jwtManager := security.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiryMinutes, cfg.JWT.Issuer, cfg.JWT.Audience, cfg.JWT.DriverAudience)

	s3Storage, err := storage.NewS3Storage(ctx, cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("init object storage: %w", err)
	}

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
	driverUpdateOnlineStatusUseCase := appdriver.NewUpdateOnlineStatusUseCase(driverRepo, vehicleRepo, kycRepo)
	driverUpdatePauseStatusUseCase := appdriver.NewUpdatePauseStatusUseCase(driverRepo)
	driverHandler := handlers.NewDriverHandler(
		driverSignupUseCase,
		driverLoginUseCase,
		driverGetProfileUseCase,
		driverUpdateProfileUseCase,
		driverUpdateOnlineStatusUseCase,
		driverUpdatePauseStatusUseCase,
	)

	vehicleRegisterUseCase := appvehicle.NewRegisterUseCase(vehicleRepo)
	vehicleListUseCase := appvehicle.NewListUseCase(vehicleRepo)
	vehicleGetUseCase := appvehicle.NewGetUseCase(vehicleRepo)
	vehicleUpdateUseCase := appvehicle.NewUpdateUseCase(vehicleRepo)
	vehicleActivateUseCase := appvehicle.NewActivateUseCase(vehicleRepo, driverRepo, kycRepo)
	vehicleDeleteUseCase := appvehicle.NewDeleteUseCase(vehicleRepo)
	vehicleHandler := handlers.NewVehicleHandler(
		vehicleRegisterUseCase,
		vehicleListUseCase,
		vehicleGetUseCase,
		vehicleUpdateUseCase,
		vehicleActivateUseCase,
		vehicleDeleteUseCase,
	)

	kycRequestUploadURLUseCase := appkyc.NewRequestUploadURLUseCase(s3Storage)
	kycConfirmUploadUseCase := appkyc.NewConfirmUploadUseCase(kycRepo, s3Storage, driverRepo, vehicleRepo)
	kycGetStatusUseCase := appkyc.NewGetStatusUseCase(kycRepo, driverRepo)
	kycHandler := handlers.NewKycHandler(
		kycRequestUploadURLUseCase,
		kycConfirmUploadUseCase,
		kycGetStatusUseCase,
	)

	placesProvider := maps.NewGoogleClient(cfg.Maps.ServerAPIKey)
	placesAutocompleteUseCase := appplaces.NewAutocompleteUseCase(placesProvider)
	placesDetailsUseCase := appplaces.NewPlaceDetailsUseCase(placesProvider)
	placesReverseGeocodeUseCase := appplaces.NewReverseGeocodeUseCase(placesProvider)
	placesHandler := handlers.NewPlacesHandler(
		placesAutocompleteUseCase,
		placesDetailsUseCase,
		placesReverseGeocodeUseCase,
	)

	savedPlaceCreateUseCase := appsavedplace.NewCreateUseCase(savedPlaceRepo)
	savedPlaceListUseCase := appsavedplace.NewListUseCase(savedPlaceRepo)
	savedPlaceUpdateUseCase := appsavedplace.NewUpdateUseCase(savedPlaceRepo)
	savedPlaceDeleteUseCase := appsavedplace.NewDeleteUseCase(savedPlaceRepo)
	savedPlaceHandler := handlers.NewSavedPlaceHandler(
		savedPlaceCreateUseCase,
		savedPlaceListUseCase,
		savedPlaceUpdateUseCase,
		savedPlaceDeleteUseCase,
	)

	router := routes.NewRouter(authHandler, driverHandler, vehicleHandler, kycHandler, placesHandler, savedPlaceHandler, jwtManager)
	return &App{Router: router}, nil
}

package routes

import (
	"net/http"

	"go-ride-backend/infrastructure/security"
	"go-ride-backend/interfaces/http/handlers"
	"go-ride-backend/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter builds the full API router. Route wiring for each domain lives
// in its own file in this package (auth_routes.go, driver_routes.go,
// vehicle_routes.go, kyc_routes.go, places_routes.go, saved_place_routes.go);
// this function just composes them.
func NewRouter(authHandler *handlers.AuthHandler, driverHandler *handlers.DriverHandler, vehicleHandler *handlers.VehicleHandler, kycHandler *handlers.KycHandler, placesHandler *handlers.PlacesHandler, savedPlaceHandler *handlers.SavedPlaceHandler, jwtManager *security.JWTManager) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := router.Group("/api/v1")
	protected := registerAuthRoutes(v1, authHandler, jwtManager)
	registerSavedPlaceRoutes(protected, savedPlaceHandler)
	driverProtected := registerDriverRoutes(v1, driverHandler, jwtManager)
	registerVehicleRoutes(driverProtected, vehicleHandler)
	registerKycRoutes(driverProtected, kycHandler)
	registerPlacesRoutes(v1, placesHandler, jwtManager)

	return router
}

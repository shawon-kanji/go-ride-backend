package routes

import (
	"go-ride-backend/infrastructure/security"
	"go-ride-backend/interfaces/http/handlers"
	"go-ride-backend/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// registerPlacesRoutes mounts the places/geocoding proxy under the generic
// authenticated group (middleware.AuthRequired, no role check) — both the
// rider and driver apps call these, unlike the role-scoped groups the other
// route files mount under.
func registerPlacesRoutes(v1 *gin.RouterGroup, placesHandler *handlers.PlacesHandler, jwtManager *security.JWTManager) {
	places := v1.Group("/places")
	places.Use(middleware.AuthRequired(jwtManager))
	{
		places.GET("/autocomplete", placesHandler.Autocomplete)
		places.GET("/reverse-geocode", placesHandler.ReverseGeocode)
		places.GET("/:place_id", placesHandler.PlaceDetails)
	}
}

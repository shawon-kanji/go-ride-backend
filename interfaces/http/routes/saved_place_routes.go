package routes

import (
	"go-ride-backend/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

// registerSavedPlaceRoutes mounts /saved-places under the rider-protected group.
func registerSavedPlaceRoutes(protected *gin.RouterGroup, savedPlaceHandler *handlers.SavedPlaceHandler) {
	savedPlaces := protected.Group("/saved-places")
	{
		savedPlaces.POST("", savedPlaceHandler.Create)
		savedPlaces.GET("", savedPlaceHandler.List)
		savedPlaces.PATCH("/:id", savedPlaceHandler.Update)
		savedPlaces.DELETE("/:id", savedPlaceHandler.Delete)
	}
}

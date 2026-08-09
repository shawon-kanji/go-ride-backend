package routes

import (
	"go-ride-backend/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

// registerVehicleRoutes mounts /vehicles under the driver-protected group.
func registerVehicleRoutes(driverProtected *gin.RouterGroup, vehicleHandler *handlers.VehicleHandler) {
	vehicles := driverProtected.Group("/vehicles")
	{
		vehicles.POST("", vehicleHandler.Register)
		vehicles.GET("", vehicleHandler.List)
		vehicles.GET("/:id", vehicleHandler.Get)
		vehicles.PATCH("/:id", vehicleHandler.Update)
		vehicles.POST("/:id/activate", vehicleHandler.Activate)
		vehicles.DELETE("/:id", vehicleHandler.Delete)
	}
}

package routes

import (
	"net/http"

	"go-ride-backend/infrastructure/security"
	"go-ride-backend/interfaces/http/handlers"
	"go-ride-backend/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *handlers.AuthHandler, driverHandler *handlers.DriverHandler, vehicleHandler *handlers.VehicleHandler, jwtManager *security.JWTManager) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.Recovery())

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
		}

		driverAuth := v1.Group("/driver/auth")
		{
			driverAuth.POST("/signup", driverHandler.Signup)
			driverAuth.POST("/login", driverHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthRequired(jwtManager))
		protected.GET("/me", authHandler.Me)
		protected.PATCH("/profile", authHandler.UpdateProfile)
		protected.POST("/change-password", authHandler.ChangePassword)
		protected.POST("/deactivate", authHandler.DeactivateAccount)

		driverProtected := v1.Group("/driver")
		driverProtected.Use(middleware.AuthRequiredRole(jwtManager, "driver"))
		driverProtected.GET("/me", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"driver_id":    c.GetString("user_id"),
				"driver_email": c.GetString("user_email"),
				"driver_role":  c.GetString("user_role"),
			})
		})
		driverProtected.GET("/profile", driverHandler.Me)
		driverProtected.PATCH("/profile", driverHandler.UpdateProfile)

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

	return router
}

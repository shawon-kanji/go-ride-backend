package routes

import (
	"net/http"

	"go-ride-backend/infrastructure/security"
	"go-ride-backend/interfaces/http/handlers"
	"go-ride-backend/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// registerDriverRoutes mounts driver auth (signup/login) and the
// driver-protected profile/online routes. It returns the driverProtected
// group so other driver-scoped route files (vehicles, kyc) can mount under
// it without re-creating the AuthRequiredRole("driver") middleware.
func registerDriverRoutes(v1 *gin.RouterGroup, driverHandler *handlers.DriverHandler, jwtManager *security.JWTManager) *gin.RouterGroup {
	driverAuth := v1.Group("/driver/auth")
	{
		driverAuth.POST("/signup", driverHandler.Signup)
		driverAuth.POST("/login", driverHandler.Login)
	}

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
	driverProtected.PATCH("/online", driverHandler.UpdateOnlineStatus)
	driverProtected.PATCH("/pause", driverHandler.UpdatePauseStatus)

	return driverProtected
}

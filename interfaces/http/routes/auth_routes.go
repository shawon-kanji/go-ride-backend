package routes

import (
	"go-ride-backend/infrastructure/security"
	"go-ride-backend/interfaces/http/handlers"
	"go-ride-backend/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// registerAuthRoutes mounts rider auth (signup/login) and the rider-protected
// profile routes (me/profile/change-password/deactivate). It returns the
// protected group so other rider-scoped route files (saved places) can mount
// under it without re-creating the AuthRequired middleware.
func registerAuthRoutes(v1 *gin.RouterGroup, authHandler *handlers.AuthHandler, jwtManager *security.JWTManager) *gin.RouterGroup {
	auth := v1.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/login", authHandler.Login)
	}

	protected := v1.Group("")
	protected.Use(middleware.AuthRequired(jwtManager))
	protected.GET("/me", authHandler.Me)
	protected.PATCH("/profile", authHandler.UpdateProfile)
	protected.POST("/change-password", authHandler.ChangePassword)
	protected.POST("/deactivate", authHandler.DeactivateAccount)

	return protected
}

package middleware

import (
	"net/http"
	"strings"

	"go-ride-backend/infrastructure/security"

	"github.com/gin-gonic/gin"
)

func AuthRequired(jwtManager *security.JWTManager) gin.HandlerFunc {
	return authRequired(jwtManager, "")
}

func AuthRequiredRole(jwtManager *security.JWTManager, requiredRole string) gin.HandlerFunc {
	return authRequired(jwtManager, requiredRole)
}

func authRequired(jwtManager *security.JWTManager, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "invalid authorization header"})
			return
		}

		claims, err := jwtManager.Parse(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "invalid or expired token"})
			return
		}

		if requiredRole != "" && claims.Role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "insufficient role"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

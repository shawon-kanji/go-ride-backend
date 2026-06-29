package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-ride-backend/infrastructure/security"

	"github.com/gin-gonic/gin"
)

func TestAuthRequiredValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := security.NewJWTManager("secret", 60, "go-ride-backend", "go-ride-clients")
	token, err := jwtManager.Generate("user-id", "user@example.com", "rider")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	r := gin.New()
	r.Use(AuthRequired(jwtManager))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAuthRequiredExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := security.NewJWTManager("secret", -1, "go-ride-backend", "go-ride-clients")
	token, err := jwtManager.Generate("user-id", "user@example.com", "rider")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	r := gin.New()
	r.Use(AuthRequired(jwtManager))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequiredMalformedToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := security.NewJWTManager("secret", 60, "go-ride-backend", "go-ride-clients")

	r := gin.New()
	r.Use(AuthRequired(jwtManager))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequiredRoleForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := security.NewJWTManager("secret", 60, "go-ride-backend", "go-ride-clients")
	token, err := jwtManager.Generate("user-id", "user@example.com", "rider")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	r := gin.New()
	r.Use(AuthRequiredRole(jwtManager, "driver"))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

package handlers

import (
	"errors"
	"net/http"

	appuser "go-ride-backend/application/user"
	"go-ride-backend/domain/user"
	"go-ride-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	signupUseCase *appuser.SignupUseCase
	loginUseCase  *appuser.LoginUseCase
}

func NewAuthHandler(signupUseCase *appuser.SignupUseCase, loginUseCase *appuser.LoginUseCase) *AuthHandler {
	return &AuthHandler{signupUseCase: signupUseCase, loginUseCase: loginUseCase}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req appuser.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	response, err := h.signupUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req appuser.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	response, err := h.loginUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	if errors.Is(err, user.ErrEmailAlreadyTaken) || errors.Is(err, user.ErrInvalidCredential) {
		httpErr := apperror.Map(err)
		c.JSON(httpErr.Status, gin.H{"code": httpErr.Code, "message": httpErr.Message})
		return
	}

	httpErr := apperror.Map(err)
	c.JSON(httpErr.Status, gin.H{"code": httpErr.Code, "message": httpErr.Message})
}

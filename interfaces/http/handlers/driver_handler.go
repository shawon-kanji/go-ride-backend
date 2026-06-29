package handlers

import (
	"net/http"

	appdriver "go-ride-backend/application/driver"
	"go-ride-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type DriverHandler struct {
	signupUseCase *appdriver.SignupUseCase
	loginUseCase  *appdriver.LoginUseCase
}

func NewDriverHandler(signupUseCase *appdriver.SignupUseCase, loginUseCase *appdriver.LoginUseCase) *DriverHandler {
	return &DriverHandler{signupUseCase: signupUseCase, loginUseCase: loginUseCase}
}

func (h *DriverHandler) Signup(c *gin.Context) {
	var req appdriver.SignupRequest
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

func (h *DriverHandler) Login(c *gin.Context) {
	var req appdriver.LoginRequest
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

func (h *DriverHandler) handleError(c *gin.Context, err error) {
	httpErr := apperror.Map(err)
	c.JSON(httpErr.Status, gin.H{"code": httpErr.Code, "message": httpErr.Message})
}

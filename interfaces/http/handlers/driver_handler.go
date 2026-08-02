package handlers

import (
	"net/http"

	appdriver "go-ride-backend/application/driver"
	"go-ride-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type DriverHandler struct {
	signupUseCase             *appdriver.SignupUseCase
	loginUseCase              *appdriver.LoginUseCase
	getProfileUseCase         *appdriver.GetProfileUseCase
	updateProfileUseCase      *appdriver.UpdateProfileUseCase
	updateOnlineStatusUseCase *appdriver.UpdateOnlineStatusUseCase
}

func NewDriverHandler(
	signupUseCase *appdriver.SignupUseCase,
	loginUseCase *appdriver.LoginUseCase,
	getProfileUseCase *appdriver.GetProfileUseCase,
	updateProfileUseCase *appdriver.UpdateProfileUseCase,
	updateOnlineStatusUseCase *appdriver.UpdateOnlineStatusUseCase,
) *DriverHandler {
	return &DriverHandler{
		signupUseCase:             signupUseCase,
		loginUseCase:              loginUseCase,
		getProfileUseCase:         getProfileUseCase,
		updateProfileUseCase:      updateProfileUseCase,
		updateOnlineStatusUseCase: updateOnlineStatusUseCase,
	}
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

func (h *DriverHandler) Me(c *gin.Context) {
	driverID := c.GetString("user_id")
	response, err := h.getProfileUseCase.Execute(c.Request.Context(), driverID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"driver": response})
}

func (h *DriverHandler) UpdateProfile(c *gin.Context) {
	var req appdriver.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	driverID := c.GetString("user_id")
	response, err := h.updateProfileUseCase.Execute(c.Request.Context(), driverID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"driver": response})
}

func (h *DriverHandler) handleError(c *gin.Context, err error) {
	httpErr := apperror.Map(err)
	c.JSON(httpErr.Status, gin.H{"code": httpErr.Code, "message": httpErr.Message})
}

func (h *DriverHandler) UpdateOnlineStatus(c *gin.Context) {
	var req appdriver.UpdateOnlineStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	driverID := c.GetString("user_id")
	response, err := h.updateOnlineStatusUseCase.Execute(c.Request.Context(), driverID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"driver": response})
}

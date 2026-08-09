package handlers

import (
	"net/http"

	appkyc "go-ride-backend/application/kyc"
	"go-ride-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type KycHandler struct {
	requestUploadURLUseCase *appkyc.RequestUploadURLUseCase
	confirmUploadUseCase    *appkyc.ConfirmUploadUseCase
	getStatusUseCase        *appkyc.GetStatusUseCase
}

func NewKycHandler(
	requestUploadURLUseCase *appkyc.RequestUploadURLUseCase,
	confirmUploadUseCase *appkyc.ConfirmUploadUseCase,
	getStatusUseCase *appkyc.GetStatusUseCase,
) *KycHandler {
	return &KycHandler{
		requestUploadURLUseCase: requestUploadURLUseCase,
		confirmUploadUseCase:    confirmUploadUseCase,
		getStatusUseCase:        getStatusUseCase,
	}
}

func (h *KycHandler) RequestUploadURL(c *gin.Context) {
	var req appkyc.RequestUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	driverID := c.GetString("user_id")
	response, err := h.requestUploadURLUseCase.Execute(c.Request.Context(), driverID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *KycHandler) ConfirmUpload(c *gin.Context) {
	var req appkyc.ConfirmUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	driverID := c.GetString("user_id")
	response, err := h.confirmUploadUseCase.Execute(c.Request.Context(), driverID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"document": response})
}

func (h *KycHandler) GetStatus(c *gin.Context) {
	driverID := c.GetString("user_id")
	response, err := h.getStatusUseCase.Execute(c.Request.Context(), driverID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *KycHandler) handleError(c *gin.Context, err error) {
	httpErr := apperror.Map(err)
	c.JSON(httpErr.Status, gin.H{"code": httpErr.Code, "message": httpErr.Message})
}

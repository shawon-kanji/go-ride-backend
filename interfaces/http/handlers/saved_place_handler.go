package handlers

import (
	"net/http"

	appsavedplace "go-ride-backend/application/savedplace"
	"go-ride-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type SavedPlaceHandler struct {
	createUseCase *appsavedplace.CreateUseCase
	listUseCase   *appsavedplace.ListUseCase
	updateUseCase *appsavedplace.UpdateUseCase
	deleteUseCase *appsavedplace.DeleteUseCase
}

func NewSavedPlaceHandler(
	createUseCase *appsavedplace.CreateUseCase,
	listUseCase *appsavedplace.ListUseCase,
	updateUseCase *appsavedplace.UpdateUseCase,
	deleteUseCase *appsavedplace.DeleteUseCase,
) *SavedPlaceHandler {
	return &SavedPlaceHandler{
		createUseCase: createUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

func (h *SavedPlaceHandler) Create(c *gin.Context) {
	var req appsavedplace.CreateSavedPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	userID := c.GetString("user_id")
	response, err := h.createUseCase.Execute(c.Request.Context(), userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"saved_place": response})
}

func (h *SavedPlaceHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	response, err := h.listUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SavedPlaceHandler) Update(c *gin.Context) {
	var req appsavedplace.UpdateSavedPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	userID := c.GetString("user_id")
	placeID := c.Param("id")
	response, err := h.updateUseCase.Execute(c.Request.Context(), userID, placeID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"saved_place": response})
}

func (h *SavedPlaceHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	placeID := c.Param("id")
	if err := h.deleteUseCase.Execute(c.Request.Context(), userID, placeID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SavedPlaceHandler) handleError(c *gin.Context, err error) {
	httpErr := apperror.Map(err)
	c.JSON(httpErr.Status, gin.H{"code": httpErr.Code, "message": httpErr.Message})
}

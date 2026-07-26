package handlers

import (
	"net/http"

	appvehicle "go-ride-backend/application/vehicle"
	"go-ride-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	registerUseCase *appvehicle.RegisterUseCase
	listUseCase     *appvehicle.ListUseCase
	getUseCase      *appvehicle.GetUseCase
	updateUseCase   *appvehicle.UpdateUseCase
	activateUseCase *appvehicle.ActivateUseCase
	deleteUseCase   *appvehicle.DeleteUseCase
}

func NewVehicleHandler(
	registerUseCase *appvehicle.RegisterUseCase,
	listUseCase *appvehicle.ListUseCase,
	getUseCase *appvehicle.GetUseCase,
	updateUseCase *appvehicle.UpdateUseCase,
	activateUseCase *appvehicle.ActivateUseCase,
	deleteUseCase *appvehicle.DeleteUseCase,
) *VehicleHandler {
	return &VehicleHandler{
		registerUseCase: registerUseCase,
		listUseCase:     listUseCase,
		getUseCase:      getUseCase,
		updateUseCase:   updateUseCase,
		activateUseCase: activateUseCase,
		deleteUseCase:   deleteUseCase,
	}
}

func (h *VehicleHandler) Register(c *gin.Context) {
	var req appvehicle.RegisterVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	driverID := c.GetString("user_id")
	response, err := h.registerUseCase.Execute(c.Request.Context(), driverID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"vehicle": response})
}

func (h *VehicleHandler) List(c *gin.Context) {
	driverID := c.GetString("user_id")
	response, err := h.listUseCase.Execute(c.Request.Context(), driverID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *VehicleHandler) Get(c *gin.Context) {
	driverID := c.GetString("user_id")
	vehicleID := c.Param("id")
	response, err := h.getUseCase.Execute(c.Request.Context(), driverID, vehicleID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"vehicle": response})
}

func (h *VehicleHandler) Update(c *gin.Context) {
	var req appvehicle.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid request payload"})
		return
	}

	driverID := c.GetString("user_id")
	vehicleID := c.Param("id")
	response, err := h.updateUseCase.Execute(c.Request.Context(), driverID, vehicleID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"vehicle": response})
}

func (h *VehicleHandler) Activate(c *gin.Context) {
	driverID := c.GetString("user_id")
	vehicleID := c.Param("id")
	response, err := h.activateUseCase.Execute(c.Request.Context(), driverID, vehicleID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"vehicle": response})
}

func (h *VehicleHandler) Delete(c *gin.Context) {
	driverID := c.GetString("user_id")
	vehicleID := c.Param("id")
	if err := h.deleteUseCase.Execute(c.Request.Context(), driverID, vehicleID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *VehicleHandler) handleError(c *gin.Context, err error) {
	httpErr := apperror.Map(err)
	c.JSON(httpErr.Status, gin.H{"code": httpErr.Code, "message": httpErr.Message})
}

package handlers

import (
	"net/http"

	appplaces "go-ride-backend/application/places"
	"go-ride-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type PlacesHandler struct {
	autocompleteUseCase   *appplaces.AutocompleteUseCase
	placeDetailsUseCase   *appplaces.PlaceDetailsUseCase
	reverseGeocodeUseCase *appplaces.ReverseGeocodeUseCase
}

func NewPlacesHandler(
	autocompleteUseCase *appplaces.AutocompleteUseCase,
	placeDetailsUseCase *appplaces.PlaceDetailsUseCase,
	reverseGeocodeUseCase *appplaces.ReverseGeocodeUseCase,
) *PlacesHandler {
	return &PlacesHandler{
		autocompleteUseCase:   autocompleteUseCase,
		placeDetailsUseCase:   placeDetailsUseCase,
		reverseGeocodeUseCase: reverseGeocodeUseCase,
	}
}

func (h *PlacesHandler) Autocomplete(c *gin.Context) {
	var req appplaces.AutocompleteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid query parameters"})
		return
	}

	response, err := h.autocompleteUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PlacesHandler) PlaceDetails(c *gin.Context) {
	req := appplaces.PlaceDetailsRequest{PlaceID: c.Param("place_id")}

	response, err := h.placeDetailsUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PlacesHandler) ReverseGeocode(c *gin.Context) {
	var req appplaces.ReverseGeocodeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": "invalid query parameters"})
		return
	}

	response, err := h.reverseGeocodeUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PlacesHandler) handleError(c *gin.Context, err error) {
	httpErr := apperror.Map(err)
	c.JSON(httpErr.Status, gin.H{"code": httpErr.Code, "message": httpErr.Message})
}

package savedplace

type CreateSavedPlaceRequest struct {
	Label            string  `json:"label" validate:"required,min=1,max=100"`
	PlaceType        string  `json:"place_type" validate:"required,oneof=home work other"`
	FormattedAddress string  `json:"formatted_address" validate:"required,min=1,max=500"`
	Lat              float64 `json:"lat" validate:"required,latitude"`
	Lng              float64 `json:"lng" validate:"required,longitude"`
}

type UpdateSavedPlaceRequest struct {
	Label            string  `json:"label" validate:"required,min=1,max=100"`
	PlaceType        string  `json:"place_type" validate:"required,oneof=home work other"`
	FormattedAddress string  `json:"formatted_address" validate:"required,min=1,max=500"`
	Lat              float64 `json:"lat" validate:"required,latitude"`
	Lng              float64 `json:"lng" validate:"required,longitude"`
}

type SavedPlaceResponse struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	Label            string  `json:"label"`
	PlaceType        string  `json:"place_type"`
	FormattedAddress string  `json:"formatted_address"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
}

type ListSavedPlacesResponse struct {
	SavedPlaces []SavedPlaceResponse `json:"saved_places"`
}

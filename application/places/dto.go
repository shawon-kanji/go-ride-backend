package places

type AutocompleteRequest struct {
	Input string   `form:"input" validate:"required,min=1,max=200"`
	Lat   *float64 `form:"lat" validate:"omitempty,min=-90,max=90"`
	Lng   *float64 `form:"lng" validate:"omitempty,min=-180,max=180"`
}

type SuggestionResponse struct {
	PlaceID       string `json:"place_id"`
	PrimaryText   string `json:"primary_text"`
	SecondaryText string `json:"secondary_text,omitempty"`
}

type AutocompleteResponse struct {
	Suggestions []SuggestionResponse `json:"suggestions"`
}

type PlaceDetailsRequest struct {
	PlaceID string `validate:"required"`
}

// Lat/Lng are pointers, not float64: go-playground validator's "required"
// tag fails on a bare float64's zero value, which would wrongly reject a
// legitimate 0.0 coordinate (equator/prime meridian) as "missing".
type ReverseGeocodeRequest struct {
	Lat *float64 `form:"lat" validate:"required,min=-90,max=90"`
	Lng *float64 `form:"lng" validate:"required,min=-180,max=180"`
}

type ResolvedPlaceResponse struct {
	PlaceID          string  `json:"place_id,omitempty"`
	Name             string  `json:"name,omitempty"`
	FormattedAddress string  `json:"formatted_address"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
}

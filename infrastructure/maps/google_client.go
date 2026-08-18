package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	appplaces "go-ride-backend/application/places"
	domainplaces "go-ride-backend/domain/places"
)

const (
	autocompleteURL = "https://maps.googleapis.com/maps/api/place/autocomplete/json"
	placeDetailsURL = "https://maps.googleapis.com/maps/api/place/details/json"
	geocodeURL      = "https://maps.googleapis.com/maps/api/geocode/json"
)

// GoogleClient implements appplaces.Provider against the legacy Places
// Autocomplete/Details and Geocoding REST APIs (not the newer places.
// googleapis.com JSON API) — simpler request shape (GET + query params vs
// POST + field mask headers) for the same data, and both are still fully
// supported on the server key provisioned for this project.
type GoogleClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewGoogleClient(apiKey string) *GoogleClient {
	return &GoogleClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type autocompleteAPIResponse struct {
	Status      string `json:"status"`
	ErrorMsg    string `json:"error_message"`
	Predictions []struct {
		PlaceID              string `json:"place_id"`
		Description          string `json:"description"`
		StructuredFormatting struct {
			MainText      string `json:"main_text"`
			SecondaryText string `json:"secondary_text"`
		} `json:"structured_formatting"`
	} `json:"predictions"`
}

func (c *GoogleClient) Autocomplete(ctx context.Context, input string, bias *appplaces.LatLng) ([]domainplaces.Suggestion, error) {
	if c.apiKey == "" {
		return nil, domainplaces.ErrProviderUnavailable
	}

	q := url.Values{}
	q.Set("input", input)
	q.Set("key", c.apiKey)
	if bias != nil {
		q.Set("location", fmt.Sprintf("%f,%f", bias.Lat, bias.Lng))
		q.Set("radius", "30000")
	}

	var body autocompleteAPIResponse
	if err := c.get(ctx, autocompleteURL, q, &body); err != nil {
		return nil, err
	}
	if err := statusToError(body.Status, body.ErrorMsg); err != nil {
		return nil, err
	}

	suggestions := make([]domainplaces.Suggestion, 0, len(body.Predictions))
	for _, p := range body.Predictions {
		primary := p.StructuredFormatting.MainText
		if primary == "" {
			primary = p.Description
		}
		suggestions = append(suggestions, domainplaces.Suggestion{
			PlaceID:       p.PlaceID,
			PrimaryText:   primary,
			SecondaryText: p.StructuredFormatting.SecondaryText,
		})
	}
	return suggestions, nil
}

type placeDetailsAPIResponse struct {
	Status   string `json:"status"`
	ErrorMsg string `json:"error_message"`
	Result   struct {
		PlaceID          string `json:"place_id"`
		Name             string `json:"name"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"result"`
}

func (c *GoogleClient) PlaceDetails(ctx context.Context, placeID string) (*domainplaces.ResolvedPlace, error) {
	if c.apiKey == "" {
		return nil, domainplaces.ErrProviderUnavailable
	}

	q := url.Values{}
	q.Set("place_id", placeID)
	q.Set("fields", "place_id,name,formatted_address,geometry")
	q.Set("key", c.apiKey)

	var body placeDetailsAPIResponse
	if err := c.get(ctx, placeDetailsURL, q, &body); err != nil {
		return nil, err
	}
	if err := statusToError(body.Status, body.ErrorMsg); err != nil {
		return nil, err
	}

	return &domainplaces.ResolvedPlace{
		PlaceID:          body.Result.PlaceID,
		Name:             body.Result.Name,
		FormattedAddress: body.Result.FormattedAddress,
		Lat:              body.Result.Geometry.Location.Lat,
		Lng:              body.Result.Geometry.Location.Lng,
	}, nil
}

type geocodeAPIResponse struct {
	Status   string `json:"status"`
	ErrorMsg string `json:"error_message"`
	Results  []struct {
		PlaceID          string `json:"place_id"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func (c *GoogleClient) ReverseGeocode(ctx context.Context, lat, lng float64) (*domainplaces.ResolvedPlace, error) {
	if c.apiKey == "" {
		return nil, domainplaces.ErrProviderUnavailable
	}

	q := url.Values{}
	q.Set("latlng", strconv.FormatFloat(lat, 'f', -1, 64)+","+strconv.FormatFloat(lng, 'f', -1, 64))
	q.Set("key", c.apiKey)

	var body geocodeAPIResponse
	if err := c.get(ctx, geocodeURL, q, &body); err != nil {
		return nil, err
	}
	// Geocoding's "no results at these coordinates" is ZERO_RESULTS with an
	// empty Results slice, not an error status — statusToError already
	// treats ZERO_RESULTS as domainplaces.ErrPlaceNotFound, so by the time
	// we reach here Results is guaranteed non-empty.
	if err := statusToError(body.Status, body.ErrorMsg); err != nil {
		return nil, err
	}

	first := body.Results[0]
	return &domainplaces.ResolvedPlace{
		PlaceID:          first.PlaceID,
		FormattedAddress: first.FormattedAddress,
		Lat:              first.Geometry.Location.Lat,
		Lng:              first.Geometry.Location.Lng,
	}, nil
}

func (c *GoogleClient) get(ctx context.Context, baseURL string, query url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"?"+query.Encode(), nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", domainplaces.ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: unexpected HTTP status %d", domainplaces.ErrProviderUnavailable, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("%w: decode response: %v", domainplaces.ErrProviderUnavailable, err)
	}
	return nil
}

// statusToError maps Google's API-level "status" field — distinct from the
// HTTP status code, which is 200 even for most of these — to our own
// sentinel errors. ZERO_RESULTS (autocomplete/geocoding) and NOT_FOUND
// (place details: a well-formed place_id the Places database doesn't have)
// both mean "this specific place didn't resolve" from our callers'
// perspective. INVALID_REQUEST is different in kind — a malformed place_id
// or parameter, i.e. bad client input, not "provider is down" or "place
// doesn't exist" — confirmed against Google's own documented status codes
// and live-tested against a garbage place_id.
func statusToError(status, errorMsg string) error {
	switch status {
	case "OK":
		return nil
	case "ZERO_RESULTS", "NOT_FOUND":
		return domainplaces.ErrPlaceNotFound
	case "INVALID_REQUEST":
		return domainplaces.ErrInvalidPlaceID
	default:
		if errorMsg != "" {
			return fmt.Errorf("%w: %s: %s", domainplaces.ErrProviderUnavailable, status, errorMsg)
		}
		return fmt.Errorf("%w: %s", domainplaces.ErrProviderUnavailable, status)
	}
}

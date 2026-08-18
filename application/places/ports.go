package places

import (
	"context"

	domainplaces "go-ride-backend/domain/places"
)

// Provider abstracts the external places/geocoding API these use cases
// need, so they depend on an interface rather than infrastructure/maps
// directly — matches application/kyc/ports.go's StoragePort precedent for
// an external, non-persisted dependency.
type Provider interface {
	Autocomplete(ctx context.Context, input string, bias *LatLng) ([]domainplaces.Suggestion, error)
	PlaceDetails(ctx context.Context, placeID string) (*domainplaces.ResolvedPlace, error)
	ReverseGeocode(ctx context.Context, lat, lng float64) (*domainplaces.ResolvedPlace, error)
}

// LatLng is the narrow coordinate shape ports.go's Provider needs — not
// domainplaces.ResolvedPlace, which carries fields (PlaceID, formatted
// address) a search bias point never has.
type LatLng struct {
	Lat float64
	Lng float64
}

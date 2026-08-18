package places

import (
	"context"
	"errors"
	"testing"

	domainplaces "go-ride-backend/domain/places"
)

// fakeProvider is defined once here and reused by place_details_test.go and
// reverse_geocode_test.go (Article XI).
type fakeProvider struct {
	autocompleteFn   func(ctx context.Context, input string, bias *LatLng) ([]domainplaces.Suggestion, error)
	placeDetailsFn   func(ctx context.Context, placeID string) (*domainplaces.ResolvedPlace, error)
	reverseGeocodeFn func(ctx context.Context, lat, lng float64) (*domainplaces.ResolvedPlace, error)
}

func (f *fakeProvider) Autocomplete(ctx context.Context, input string, bias *LatLng) ([]domainplaces.Suggestion, error) {
	return f.autocompleteFn(ctx, input, bias)
}

func (f *fakeProvider) PlaceDetails(ctx context.Context, placeID string) (*domainplaces.ResolvedPlace, error) {
	return f.placeDetailsFn(ctx, placeID)
}

func (f *fakeProvider) ReverseGeocode(ctx context.Context, lat, lng float64) (*domainplaces.ResolvedPlace, error) {
	return f.reverseGeocodeFn(ctx, lat, lng)
}

func TestAutocompleteUseCaseExecuteSuccess(t *testing.T) {
	provider := &fakeProvider{
		autocompleteFn: func(ctx context.Context, input string, bias *LatLng) ([]domainplaces.Suggestion, error) {
			if input != "Orchard" {
				t.Fatalf("expected input %q, got %q", "Orchard", input)
			}
			if bias == nil || bias.Lat != 1.35 || bias.Lng != 103.8 {
				t.Fatalf("expected bias {1.35 103.8}, got %+v", bias)
			}
			return []domainplaces.Suggestion{{PlaceID: "p1", PrimaryText: "Orchard Road", SecondaryText: "Singapore"}}, nil
		},
	}
	uc := NewAutocompleteUseCase(provider)

	lat, lng := 1.35, 103.8
	resp, err := uc.Execute(context.Background(), AutocompleteRequest{Input: "Orchard", Lat: &lat, Lng: &lng})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Suggestions) != 1 || resp.Suggestions[0].PlaceID != "p1" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestAutocompleteUseCaseExecuteEmptyInputRejected(t *testing.T) {
	provider := &fakeProvider{
		autocompleteFn: func(ctx context.Context, input string, bias *LatLng) ([]domainplaces.Suggestion, error) {
			t.Fatal("provider should not be called for an invalid request")
			return nil, nil
		},
	}
	uc := NewAutocompleteUseCase(provider)

	_, err := uc.Execute(context.Background(), AutocompleteRequest{Input: ""})
	var verr ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestAutocompleteUseCaseExecuteProviderError(t *testing.T) {
	provider := &fakeProvider{
		autocompleteFn: func(ctx context.Context, input string, bias *LatLng) ([]domainplaces.Suggestion, error) {
			return nil, domainplaces.ErrProviderUnavailable
		},
	}
	uc := NewAutocompleteUseCase(provider)

	_, err := uc.Execute(context.Background(), AutocompleteRequest{Input: "Orchard"})
	if !errors.Is(err, domainplaces.ErrProviderUnavailable) {
		t.Fatalf("expected ErrProviderUnavailable, got %v", err)
	}
}

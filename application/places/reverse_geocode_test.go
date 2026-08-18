package places

import (
	"context"
	"errors"
	"testing"

	domainplaces "go-ride-backend/domain/places"
)

func TestReverseGeocodeUseCaseExecuteSuccess(t *testing.T) {
	provider := &fakeProvider{
		reverseGeocodeFn: func(ctx context.Context, lat, lng float64) (*domainplaces.ResolvedPlace, error) {
			if lat != 1.3 || lng != 103.8 {
				t.Fatalf("expected (1.3, 103.8), got (%v, %v)", lat, lng)
			}
			return &domainplaces.ResolvedPlace{FormattedAddress: "31a Scotts Rd, Singapore", Lat: lat, Lng: lng}, nil
		},
	}
	uc := NewReverseGeocodeUseCase(provider)

	lat, lng := 1.3, 103.8
	resp, err := uc.Execute(context.Background(), ReverseGeocodeRequest{Lat: &lat, Lng: &lng})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.FormattedAddress != "31a Scotts Rd, Singapore" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

// TestReverseGeocodeUseCaseExecuteZeroCoordinateAccepted guards the exact
// bug dto.go's comment on Lat/Lng warns about: a bare float64 "required"
// tag would wrongly reject 0.0 as "missing". Lat/Lng are pointers so a
// legitimate equator/prime-meridian coordinate passes validation.
func TestReverseGeocodeUseCaseExecuteZeroCoordinateAccepted(t *testing.T) {
	provider := &fakeProvider{
		reverseGeocodeFn: func(ctx context.Context, lat, lng float64) (*domainplaces.ResolvedPlace, error) {
			return &domainplaces.ResolvedPlace{FormattedAddress: "Gulf of Guinea", Lat: lat, Lng: lng}, nil
		},
	}
	uc := NewReverseGeocodeUseCase(provider)

	zero := 0.0
	_, err := uc.Execute(context.Background(), ReverseGeocodeRequest{Lat: &zero, Lng: &zero})
	if err != nil {
		t.Fatalf("unexpected error for a valid zero coordinate: %v", err)
	}
}

func TestReverseGeocodeUseCaseExecuteMissingCoordinatesRejected(t *testing.T) {
	provider := &fakeProvider{
		reverseGeocodeFn: func(ctx context.Context, lat, lng float64) (*domainplaces.ResolvedPlace, error) {
			t.Fatal("provider should not be called for an invalid request")
			return nil, nil
		},
	}
	uc := NewReverseGeocodeUseCase(provider)

	_, err := uc.Execute(context.Background(), ReverseGeocodeRequest{})
	var verr ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestReverseGeocodeUseCaseExecuteOutOfRangeRejected(t *testing.T) {
	provider := &fakeProvider{
		reverseGeocodeFn: func(ctx context.Context, lat, lng float64) (*domainplaces.ResolvedPlace, error) {
			t.Fatal("provider should not be called for an invalid request")
			return nil, nil
		},
	}
	uc := NewReverseGeocodeUseCase(provider)

	lat, lng := 200.0, 103.8
	_, err := uc.Execute(context.Background(), ReverseGeocodeRequest{Lat: &lat, Lng: &lng})
	var verr ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

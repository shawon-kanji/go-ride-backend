package places

import (
	"context"
	"errors"
	"testing"

	domainplaces "go-ride-backend/domain/places"
)

func TestPlaceDetailsUseCaseExecuteSuccess(t *testing.T) {
	provider := &fakeProvider{
		placeDetailsFn: func(ctx context.Context, placeID string) (*domainplaces.ResolvedPlace, error) {
			if placeID != "p1" {
				t.Fatalf("expected place_id %q, got %q", "p1", placeID)
			}
			return &domainplaces.ResolvedPlace{PlaceID: "p1", Name: "Orchard Road", FormattedAddress: "Orchard Rd, Singapore", Lat: 1.3, Lng: 103.8}, nil
		},
	}
	uc := NewPlaceDetailsUseCase(provider)

	resp, err := uc.Execute(context.Background(), PlaceDetailsRequest{PlaceID: "p1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.FormattedAddress != "Orchard Rd, Singapore" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestPlaceDetailsUseCaseExecuteEmptyPlaceIDRejected(t *testing.T) {
	provider := &fakeProvider{
		placeDetailsFn: func(ctx context.Context, placeID string) (*domainplaces.ResolvedPlace, error) {
			t.Fatal("provider should not be called for an invalid request")
			return nil, nil
		},
	}
	uc := NewPlaceDetailsUseCase(provider)

	_, err := uc.Execute(context.Background(), PlaceDetailsRequest{PlaceID: ""})
	var verr ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestPlaceDetailsUseCaseExecuteNotFound(t *testing.T) {
	provider := &fakeProvider{
		placeDetailsFn: func(ctx context.Context, placeID string) (*domainplaces.ResolvedPlace, error) {
			return nil, domainplaces.ErrPlaceNotFound
		},
	}
	uc := NewPlaceDetailsUseCase(provider)

	_, err := uc.Execute(context.Background(), PlaceDetailsRequest{PlaceID: "unknown"})
	if !errors.Is(err, domainplaces.ErrPlaceNotFound) {
		t.Fatalf("expected ErrPlaceNotFound, got %v", err)
	}
}

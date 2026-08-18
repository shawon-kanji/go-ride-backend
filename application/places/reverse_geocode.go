package places

import "context"

type ReverseGeocodeUseCase struct {
	provider Provider
}

func NewReverseGeocodeUseCase(provider Provider) *ReverseGeocodeUseCase {
	return &ReverseGeocodeUseCase{provider: provider}
}

func (uc *ReverseGeocodeUseCase) Execute(ctx context.Context, req ReverseGeocodeRequest) (*ResolvedPlaceResponse, error) {
	if err := validateReverseGeocodeRequest(req); err != nil {
		return nil, err
	}

	place, err := uc.provider.ReverseGeocode(ctx, *req.Lat, *req.Lng)
	if err != nil {
		return nil, err
	}

	resp := ToResolvedPlaceResponse(place)
	return &resp, nil
}

package places

import "context"

type PlaceDetailsUseCase struct {
	provider Provider
}

func NewPlaceDetailsUseCase(provider Provider) *PlaceDetailsUseCase {
	return &PlaceDetailsUseCase{provider: provider}
}

func (uc *PlaceDetailsUseCase) Execute(ctx context.Context, req PlaceDetailsRequest) (*ResolvedPlaceResponse, error) {
	if err := validatePlaceDetailsRequest(req); err != nil {
		return nil, err
	}

	place, err := uc.provider.PlaceDetails(ctx, req.PlaceID)
	if err != nil {
		return nil, err
	}

	resp := ToResolvedPlaceResponse(place)
	return &resp, nil
}

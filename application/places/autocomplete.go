package places

import "context"

type AutocompleteUseCase struct {
	provider Provider
}

func NewAutocompleteUseCase(provider Provider) *AutocompleteUseCase {
	return &AutocompleteUseCase{provider: provider}
}

func (uc *AutocompleteUseCase) Execute(ctx context.Context, req AutocompleteRequest) (*AutocompleteResponse, error) {
	if err := validateAutocompleteRequest(req); err != nil {
		return nil, err
	}

	var bias *LatLng
	if req.Lat != nil && req.Lng != nil {
		bias = &LatLng{Lat: *req.Lat, Lng: *req.Lng}
	}

	suggestions, err := uc.provider.Autocomplete(ctx, req.Input, bias)
	if err != nil {
		return nil, err
	}

	resp := &AutocompleteResponse{Suggestions: make([]SuggestionResponse, 0, len(suggestions))}
	for _, s := range suggestions {
		resp.Suggestions = append(resp.Suggestions, ToSuggestionResponse(s))
	}
	return resp, nil
}

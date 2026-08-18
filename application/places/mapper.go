package places

import domainplaces "go-ride-backend/domain/places"

func ToSuggestionResponse(s domainplaces.Suggestion) SuggestionResponse {
	return SuggestionResponse{
		PlaceID:       s.PlaceID,
		PrimaryText:   s.PrimaryText,
		SecondaryText: s.SecondaryText,
	}
}

func ToResolvedPlaceResponse(p *domainplaces.ResolvedPlace) ResolvedPlaceResponse {
	return ResolvedPlaceResponse{
		PlaceID:          p.PlaceID,
		Name:             p.Name,
		FormattedAddress: p.FormattedAddress,
		Lat:              p.Lat,
		Lng:              p.Lng,
	}
}

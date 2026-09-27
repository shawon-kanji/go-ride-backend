package savedplace

import (
	domainsavedplace "go-ride-backend/domain/savedplace"
)

func ToSavedPlaceResponse(p *domainsavedplace.SavedPlace) SavedPlaceResponse {
	return SavedPlaceResponse{
		ID:               p.ID.String(),
		UserID:           p.UserID.String(),
		Label:            p.Label,
		PlaceType:        p.PlaceType,
		FormattedAddress: p.FormattedAddress,
		Lat:              p.Lat,
		Lng:              p.Lng,
	}
}

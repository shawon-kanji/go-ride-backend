package savedplace

import (
	"context"

	domainsavedplace "go-ride-backend/domain/savedplace"

	"github.com/google/uuid"
)

func parseIDs(userIDStr, placeIDStr string) (userID uuid.UUID, placeID uuid.UUID, err error) {
	userID, err = uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, uuid.Nil, ValidationError{Message: "invalid user id"}
	}
	placeID, err = uuid.Parse(placeIDStr)
	if err != nil {
		return uuid.Nil, uuid.Nil, ValidationError{Message: "invalid saved place id"}
	}
	return userID, placeID, nil
}

// loadOwnedSavedPlace fetches a saved place and verifies it belongs to
// userID, returning ErrSavedPlaceForbidden rather than leaking existence of
// places owned by other users.
func loadOwnedSavedPlace(ctx context.Context, repo domainsavedplace.Repository, userID, placeID uuid.UUID) (*domainsavedplace.SavedPlace, error) {
	p, err := repo.GetByID(ctx, placeID)
	if err != nil {
		return nil, err
	}
	if p.UserID != userID {
		return nil, domainsavedplace.ErrSavedPlaceForbidden
	}
	return p, nil
}

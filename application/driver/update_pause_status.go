package driver

import (
	"context"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

type UpdatePauseStatusUseCase struct {
	repo domaindriver.Repository
}

func NewUpdatePauseStatusUseCase(repo domaindriver.Repository) *UpdatePauseStatusUseCase {
	return &UpdatePauseStatusUseCase{repo: repo}
}

// Execute rejects pausing (IsPaused: true) a driver who isn't currently
// online -- D08's Pause chip is explicitly "distinct from going offline", a
// break within a shift, not a way to start one. Unpausing has no such
// restriction: it's also how UpdateOnlineStatus's own transaction clears a
// stale pause when a driver goes offline, so an already-offline driver
// unpausing (a redundant no-op from the client's perspective) must not
// error.
func (uc *UpdatePauseStatusUseCase) Execute(ctx context.Context, driverID string, req UpdatePauseStatusRequest) (*DriverResponse, error) {
	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	if req.IsPaused {
		driver, err := uc.repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if !driver.IsOnline {
			return nil, domaindriver.ErrDriverNotOnline
		}
	}

	updated, err := uc.repo.UpdatePauseStatus(ctx, id, req.IsPaused)
	if err != nil {
		return nil, err
	}

	resp := ToDriverResponse(updated)
	return &resp, nil
}

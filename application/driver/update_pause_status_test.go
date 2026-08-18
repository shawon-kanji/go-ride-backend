package driver

import (
	"context"
	"errors"
	"testing"
	"time"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

func TestUpdatePauseStatusUseCaseExecutePauseSuccess(t *testing.T) {
	driverID := uuid.New()
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"driver@example.com": {ID: driverID, Email: "driver@example.com", IsOnline: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}}
	uc := NewUpdatePauseStatusUseCase(repo)

	resp, err := uc.Execute(context.Background(), driverID.String(), UpdatePauseStatusRequest{IsPaused: true})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !resp.IsPaused {
		t.Errorf("IsPaused = false, want true")
	}
}

func TestUpdatePauseStatusUseCaseExecuteUnpauseAlwaysAllowed(t *testing.T) {
	driverID := uuid.New()
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		// Offline driver unpausing (e.g. a redundant client resync) must not error.
		"driver@example.com": {ID: driverID, Email: "driver@example.com", IsOnline: false, IsPaused: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}}
	uc := NewUpdatePauseStatusUseCase(repo)

	resp, err := uc.Execute(context.Background(), driverID.String(), UpdatePauseStatusRequest{IsPaused: false})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if resp.IsPaused {
		t.Errorf("IsPaused = true, want false")
	}
}

func TestUpdatePauseStatusUseCaseExecuteBlockedWhenNotOnline(t *testing.T) {
	driverID := uuid.New()
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"driver@example.com": {ID: driverID, Email: "driver@example.com", IsOnline: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}}
	uc := NewUpdatePauseStatusUseCase(repo)

	_, err := uc.Execute(context.Background(), driverID.String(), UpdatePauseStatusRequest{IsPaused: true})
	if !errors.Is(err, domaindriver.ErrDriverNotOnline) {
		t.Fatalf("err = %v, want ErrDriverNotOnline", err)
	}
}

func TestUpdatePauseStatusUseCaseExecuteDriverNotFound(t *testing.T) {
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{}}
	uc := NewUpdatePauseStatusUseCase(repo)

	_, err := uc.Execute(context.Background(), uuid.New().String(), UpdatePauseStatusRequest{IsPaused: true})
	if !errors.Is(err, domaindriver.ErrDriverNotFound) {
		t.Fatalf("err = %v, want ErrDriverNotFound", err)
	}
}

func TestUpdatePauseStatusUseCaseExecuteInvalidDriverID(t *testing.T) {
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{}}
	uc := NewUpdatePauseStatusUseCase(repo)

	_, err := uc.Execute(context.Background(), "not-a-uuid", UpdatePauseStatusRequest{IsPaused: true})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("err = %v, want ValidationError", err)
	}
}

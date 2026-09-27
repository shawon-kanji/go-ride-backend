package savedplace

import (
	"context"
	"errors"
	"testing"

	domainsavedplace "go-ride-backend/domain/savedplace"

	"github.com/google/uuid"
)

func TestUpdateUseCaseExecuteSuccess(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	createUC := NewCreateUseCase(repo)
	updateUC := NewUpdateUseCase(repo)
	userID := uuid.New()

	created, err := createUC.Execute(context.Background(), userID.String(), validCreateReq("other"))
	if err != nil {
		t.Fatalf("expected no error creating place, got %v", err)
	}

	req := UpdateSavedPlaceRequest{
		Label:            "Gym",
		PlaceType:        "other",
		FormattedAddress: "2 Second St",
		Lat:              3.15,
		Lng:              101.7,
	}
	resp, err := updateUC.Execute(context.Background(), userID.String(), created.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Label != "Gym" || resp.FormattedAddress != "2 Second St" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestUpdateUseCaseExecuteForbidden(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	createUC := NewCreateUseCase(repo)
	updateUC := NewUpdateUseCase(repo)
	owner := uuid.New()
	intruder := uuid.New()

	created, err := createUC.Execute(context.Background(), owner.String(), validCreateReq("home"))
	if err != nil {
		t.Fatalf("expected no error creating place, got %v", err)
	}

	_, err = updateUC.Execute(context.Background(), intruder.String(), created.ID, UpdateSavedPlaceRequest{
		Label: "Hijacked", PlaceType: "home", FormattedAddress: "x", Lat: 1, Lng: 1,
	})
	if !errors.Is(err, domainsavedplace.ErrSavedPlaceForbidden) {
		t.Fatalf("expected ErrSavedPlaceForbidden, got %v", err)
	}
}

func TestUpdateUseCaseExecuteNotFound(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	uc := NewUpdateUseCase(repo)

	_, err := uc.Execute(context.Background(), uuid.New().String(), uuid.New().String(), UpdateSavedPlaceRequest{
		Label: "Ghost", PlaceType: "home", FormattedAddress: "x", Lat: 1, Lng: 1,
	})
	if !errors.Is(err, domainsavedplace.ErrSavedPlaceNotFound) {
		t.Fatalf("expected ErrSavedPlaceNotFound, got %v", err)
	}
}

func TestUpdateUseCaseExecuteTypeAlreadyTaken(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	createUC := NewCreateUseCase(repo)
	updateUC := NewUpdateUseCase(repo)
	userID := uuid.New()

	if _, err := createUC.Execute(context.Background(), userID.String(), validCreateReq("home")); err != nil {
		t.Fatalf("expected no error creating home, got %v", err)
	}
	work, err := createUC.Execute(context.Background(), userID.String(), validCreateReq("work"))
	if err != nil {
		t.Fatalf("expected no error creating work, got %v", err)
	}

	_, err = updateUC.Execute(context.Background(), userID.String(), work.ID, UpdateSavedPlaceRequest{
		Label: "Work", PlaceType: "home", FormattedAddress: "x", Lat: 1, Lng: 1,
	})
	if !errors.Is(err, domainsavedplace.ErrPlaceTypeTaken) {
		t.Fatalf("expected ErrPlaceTypeTaken, got %v", err)
	}
}

func TestUpdateUseCaseExecuteValidation(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	createUC := NewCreateUseCase(repo)
	updateUC := NewUpdateUseCase(repo)
	userID := uuid.New()

	created, err := createUC.Execute(context.Background(), userID.String(), validCreateReq("other"))
	if err != nil {
		t.Fatalf("expected no error creating place, got %v", err)
	}

	_, err = updateUC.Execute(context.Background(), userID.String(), created.ID, UpdateSavedPlaceRequest{
		Label: "", PlaceType: "other", FormattedAddress: "x", Lat: 1, Lng: 1,
	})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

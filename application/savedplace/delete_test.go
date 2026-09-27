package savedplace

import (
	"context"
	"errors"
	"testing"

	domainsavedplace "go-ride-backend/domain/savedplace"

	"github.com/google/uuid"
)

func TestDeleteUseCaseExecuteSuccess(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	createUC := NewCreateUseCase(repo)
	deleteUC := NewDeleteUseCase(repo)
	userID := uuid.New()

	created, err := createUC.Execute(context.Background(), userID.String(), validCreateReq("other"))
	if err != nil {
		t.Fatalf("expected no error creating place, got %v", err)
	}

	if err := deleteUC.Execute(context.Background(), userID.String(), created.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := repo.GetByID(context.Background(), uuid.MustParse(created.ID)); !errors.Is(err, domainsavedplace.ErrSavedPlaceNotFound) {
		t.Fatalf("expected place to be gone, got %v", err)
	}
}

func TestDeleteUseCaseExecuteForbidden(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	createUC := NewCreateUseCase(repo)
	deleteUC := NewDeleteUseCase(repo)
	owner := uuid.New()
	intruder := uuid.New()

	created, err := createUC.Execute(context.Background(), owner.String(), validCreateReq("home"))
	if err != nil {
		t.Fatalf("expected no error creating place, got %v", err)
	}

	err = deleteUC.Execute(context.Background(), intruder.String(), created.ID)
	if !errors.Is(err, domainsavedplace.ErrSavedPlaceForbidden) {
		t.Fatalf("expected ErrSavedPlaceForbidden, got %v", err)
	}
}

func TestDeleteUseCaseExecuteNotFound(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	uc := NewDeleteUseCase(repo)

	err := uc.Execute(context.Background(), uuid.New().String(), uuid.New().String())
	if !errors.Is(err, domainsavedplace.ErrSavedPlaceNotFound) {
		t.Fatalf("expected ErrSavedPlaceNotFound, got %v", err)
	}
}

func TestDeleteUseCaseExecuteInvalidID(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	uc := NewDeleteUseCase(repo)

	err := uc.Execute(context.Background(), uuid.New().String(), "not-a-uuid")
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

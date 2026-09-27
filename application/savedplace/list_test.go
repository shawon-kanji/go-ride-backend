package savedplace

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestListUseCaseExecuteSuccess(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	createUC := NewCreateUseCase(repo)
	listUC := NewListUseCase(repo)
	userID := uuid.New()

	if _, err := createUC.Execute(context.Background(), userID.String(), validCreateReq("home")); err != nil {
		t.Fatalf("expected no error creating place, got %v", err)
	}
	if _, err := createUC.Execute(context.Background(), userID.String(), validCreateReq("work")); err != nil {
		t.Fatalf("expected no error creating place, got %v", err)
	}
	if _, err := createUC.Execute(context.Background(), uuid.New().String(), validCreateReq("home")); err != nil {
		t.Fatalf("expected no error creating another user's place, got %v", err)
	}

	resp, err := listUC.Execute(context.Background(), userID.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.SavedPlaces) != 2 {
		t.Fatalf("expected 2 saved places for user, got %d", len(resp.SavedPlaces))
	}
}

func TestListUseCaseExecuteInvalidUserID(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	uc := NewListUseCase(repo)

	_, err := uc.Execute(context.Background(), "not-a-uuid")
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

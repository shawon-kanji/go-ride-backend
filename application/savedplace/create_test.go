package savedplace

import (
	"context"
	"errors"
	"testing"
	"time"

	domainsavedplace "go-ride-backend/domain/savedplace"

	"github.com/google/uuid"
)

type fakeSavedPlaceRepo struct {
	byID map[uuid.UUID]*domainsavedplace.SavedPlace
}

func newFakeSavedPlaceRepo() *fakeSavedPlaceRepo {
	return &fakeSavedPlaceRepo{byID: map[uuid.UUID]*domainsavedplace.SavedPlace{}}
}

func (r *fakeSavedPlaceRepo) Create(_ context.Context, input *domainsavedplace.CreateInput) (*domainsavedplace.SavedPlace, error) {
	p := &domainsavedplace.SavedPlace{
		ID:               uuid.New(),
		UserID:           input.UserID,
		Label:            input.Label,
		PlaceType:        input.PlaceType,
		FormattedAddress: input.FormattedAddress,
		Lat:              input.Lat,
		Lng:              input.Lng,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	r.byID[p.ID] = p
	return p, nil
}

func (r *fakeSavedPlaceRepo) GetByID(_ context.Context, id uuid.UUID) (*domainsavedplace.SavedPlace, error) {
	if p, ok := r.byID[id]; ok {
		return p, nil
	}
	return nil, domainsavedplace.ErrSavedPlaceNotFound
}

func (r *fakeSavedPlaceRepo) GetByUserAndType(_ context.Context, userID uuid.UUID, placeType string) (*domainsavedplace.SavedPlace, error) {
	for _, p := range r.byID {
		if p.UserID == userID && p.PlaceType == placeType {
			return p, nil
		}
	}
	return nil, domainsavedplace.ErrSavedPlaceNotFound
}

func (r *fakeSavedPlaceRepo) ListByUser(_ context.Context, userID uuid.UUID) ([]*domainsavedplace.SavedPlace, error) {
	var out []*domainsavedplace.SavedPlace
	for _, p := range r.byID {
		if p.UserID == userID {
			out = append(out, p)
		}
	}
	return out, nil
}

func (r *fakeSavedPlaceRepo) Update(_ context.Context, id uuid.UUID, input *domainsavedplace.UpdateInput) (*domainsavedplace.SavedPlace, error) {
	p, ok := r.byID[id]
	if !ok {
		return nil, domainsavedplace.ErrSavedPlaceNotFound
	}
	p.Label = input.Label
	p.PlaceType = input.PlaceType
	p.FormattedAddress = input.FormattedAddress
	p.Lat = input.Lat
	p.Lng = input.Lng
	return p, nil
}

func (r *fakeSavedPlaceRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := r.byID[id]; !ok {
		return domainsavedplace.ErrSavedPlaceNotFound
	}
	delete(r.byID, id)
	return nil
}

func validCreateReq(placeType string) CreateSavedPlaceRequest {
	return CreateSavedPlaceRequest{
		Label:            "Home",
		PlaceType:        placeType,
		FormattedAddress: "1 Main St",
		Lat:              3.139,
		Lng:              101.6869,
	}
}

func TestCreateUseCaseExecuteSuccess(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	uc := NewCreateUseCase(repo)
	userID := uuid.New()

	resp, err := uc.Execute(context.Background(), userID.String(), validCreateReq("home"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.PlaceType != "home" || resp.UserID != userID.String() {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestCreateUseCaseExecuteHomeAlreadyTaken(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	uc := NewCreateUseCase(repo)
	userID := uuid.New()

	if _, err := uc.Execute(context.Background(), userID.String(), validCreateReq("home")); err != nil {
		t.Fatalf("expected no error on first create, got %v", err)
	}

	_, err := uc.Execute(context.Background(), userID.String(), validCreateReq("home"))
	if !errors.Is(err, domainsavedplace.ErrPlaceTypeTaken) {
		t.Fatalf("expected ErrPlaceTypeTaken, got %v", err)
	}
}

func TestCreateUseCaseExecuteMultipleOtherAllowed(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	uc := NewCreateUseCase(repo)
	userID := uuid.New()

	if _, err := uc.Execute(context.Background(), userID.String(), validCreateReq("other")); err != nil {
		t.Fatalf("expected no error on first create, got %v", err)
	}
	if _, err := uc.Execute(context.Background(), userID.String(), validCreateReq("other")); err != nil {
		t.Fatalf("expected a second 'other' place to be allowed, got %v", err)
	}
}

func TestCreateUseCaseExecuteValidation(t *testing.T) {
	repo := newFakeSavedPlaceRepo()
	uc := NewCreateUseCase(repo)

	req := validCreateReq("home")
	req.Label = ""
	_, err := uc.Execute(context.Background(), uuid.New().String(), req)
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

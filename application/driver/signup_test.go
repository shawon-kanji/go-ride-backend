package driver

import (
	"context"
	"errors"
	"testing"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

type fakeDriverRepo struct {
	byEmail map[string]*domaindriver.Driver
}

func (r *fakeDriverRepo) Create(_ context.Context, driver *domaindriver.Driver) error {
	r.byEmail[driver.Email] = driver
	return nil
}

func (r *fakeDriverRepo) GetByEmail(_ context.Context, email string) (*domaindriver.Driver, error) {
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, domaindriver.ErrDriverNotFound
}

func (r *fakeDriverRepo) GetByID(_ context.Context, id uuid.UUID) (*domaindriver.Driver, error) {
	for _, d := range r.byEmail {
		if d.ID == id {
			return d, nil
		}
	}
	return nil, domaindriver.ErrDriverNotFound
}

func (r *fakeDriverRepo) UpdateProfile(_ context.Context, id uuid.UUID, firstName string, lastName string) (*domaindriver.Driver, error) {
	for _, d := range r.byEmail {
		if d.ID == id {
			d.FirstName = firstName
			d.LastName = lastName
			return d, nil
		}
	}
	return nil, domaindriver.ErrDriverNotFound
}

func (r *fakeDriverRepo) UpdateOnlineStatus(_ context.Context, id uuid.UUID, isOnline bool) (*domaindriver.Driver, error) {
	for _, d := range r.byEmail {
		if d.ID == id {
			d.IsOnline = isOnline
			return d, nil
		}
	}
	return nil, domaindriver.ErrDriverNotFound
}

type fakeHasher struct{}

func (h *fakeHasher) Hash(password string) (string, error) {
	return "hashed:" + password, nil
}

func (h *fakeHasher) Compare(hashedPassword string, password string) error {
	if hashedPassword != "hashed:"+password {
		return errors.New("invalid password")
	}
	return nil
}

type fakeTokens struct{}

func (t *fakeTokens) Generate(userID string, email string, role string) (string, error) {
	return userID + ":" + email + ":" + role, nil
}

func TestSignupUseCaseExecuteSuccess(t *testing.T) {
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{}}
	hasher := &fakeHasher{}
	uc := NewSignupUseCase(repo, hasher)

	resp, err := uc.Execute(context.Background(), SignupRequest{
		Email:     "foo@example.com",
		Password:  "password123",
		FirstName: "Foo",
		LastName:  "Bar",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.Driver.Email != "foo@example.com" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if repo.byEmail["foo@example.com"] == nil {
		t.Fatalf("driver was not persisted")
	}
	if resp.Driver.IsEmailVerified {
		t.Fatalf("expected is_email_verified=false by default")
	}
	if resp.Driver.AccountStatus != domaindriver.AccountStatusPending {
		t.Fatalf("expected account status pending, got %s", resp.Driver.AccountStatus)
	}
}

func TestSignupUseCaseExecuteDuplicateEmail(t *testing.T) {
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"foo@example.com": {
			ID:           uuid.New(),
			Email:        "foo@example.com",
			PasswordHash: "hashed:password123",
			FirstName:    "Foo",
			LastName:     "Bar",
		},
	}}
	hasher := &fakeHasher{}
	uc := NewSignupUseCase(repo, hasher)

	_, err := uc.Execute(context.Background(), SignupRequest{
		Email:     "foo@example.com",
		Password:  "password123",
		FirstName: "Foo",
		LastName:  "Bar",
	})
	if !errors.Is(err, domaindriver.ErrEmailAlreadyTaken) {
		t.Fatalf("expected ErrEmailAlreadyTaken, got %v", err)
	}
}

func TestSignupUseCaseExecuteValidation(t *testing.T) {
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{}}
	hasher := &fakeHasher{}
	uc := NewSignupUseCase(repo, hasher)

	_, err := uc.Execute(context.Background(), SignupRequest{
		Email:     "not-email",
		Password:  "short",
		FirstName: "A",
		LastName:  "B",
	})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

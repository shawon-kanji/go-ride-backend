package user

import (
	"context"
	"errors"
	"testing"
	"time"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

type fakeUserRepo struct {
	byEmail map[string]*domainuser.User
}

func (r *fakeUserRepo) Create(_ context.Context, input *domainuser.SignupInput) (*domainuser.User, error) {
	user := &domainuser.User{
		ID:            uuid.New(),
		Email:         input.Email,
		PasswordHash:  input.PasswordHash,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		AccountStatus: input.AccountStatus,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	r.byEmail[user.Email] = user
	return user, nil
}

func (r *fakeUserRepo) GetByEmail(_ context.Context, email string) (*domainuser.User, error) {
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, domainuser.ErrUserNotFound
}

func (r *fakeUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domainuser.User, error) {
	for _, u := range r.byEmail {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, domainuser.ErrUserNotFound
}

func (r *fakeUserRepo) UpdateProfile(_ context.Context, id uuid.UUID, firstName string, lastName string) (*domainuser.User, error) {
	for _, u := range r.byEmail {
		if u.ID == id {
			u.FirstName = firstName
			u.LastName = lastName
			u.UpdatedAt = time.Now().UTC()
			return u, nil
		}
	}
	return nil, domainuser.ErrUserNotFound
}

func (r *fakeUserRepo) UpdatePassword(_ context.Context, id uuid.UUID, passwordHash string) error {
	for _, u := range r.byEmail {
		if u.ID == id {
			u.PasswordHash = passwordHash
			u.UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return domainuser.ErrUserNotFound
}

func (r *fakeUserRepo) Deactivate(_ context.Context, id uuid.UUID, at time.Time) error {
	for _, u := range r.byEmail {
		if u.ID == id {
			u.AccountStatus = domainuser.AccountStatusDeactivated
			u.DeactivatedAt = &at
			u.UpdatedAt = at
			return nil
		}
	}
	return domainuser.ErrUserNotFound
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
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{}}
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
	if resp == nil || resp.User.Email != "foo@example.com" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if repo.byEmail["foo@example.com"] == nil {
		t.Fatalf("user was not persisted")
	}
}

func TestSignupUseCaseExecuteDuplicateEmail(t *testing.T) {
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{
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
	if !errors.Is(err, domainuser.ErrEmailAlreadyTaken) {
		t.Fatalf("expected ErrEmailAlreadyTaken, got %v", err)
	}
}

func TestSignupUseCaseExecuteValidation(t *testing.T) {
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{}}
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

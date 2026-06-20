package user

import (
	"context"
	"errors"
	"testing"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

type fakeUserRepo struct {
	byEmail map[string]*domainuser.User
}

func (r *fakeUserRepo) Create(_ context.Context, user *domainuser.User) error {
	r.byEmail[user.Email] = user
	return nil
}

func (r *fakeUserRepo) GetByEmail(_ context.Context, email string) (*domainuser.User, error) {
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, domainuser.ErrUserNotFound
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

func (t *fakeTokens) Generate(userID string, email string) (string, error) {
	return userID + ":" + email, nil
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

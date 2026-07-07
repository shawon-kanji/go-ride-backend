package user

import (
	"strings"

	domainuser "go-ride-backend/domain/user"
)

func NewUserFromSignup(req SignupRequest, hashedPassword string) *domainuser.SignupInput {

	return &domainuser.SignupInput{
		Email:         strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash:  hashedPassword,
		FirstName:     strings.TrimSpace(req.FirstName),
		LastName:      strings.TrimSpace(req.LastName),
		AccountStatus: domainuser.AccountStatusActive,
	}
}

func ToUserResponse(u *domainuser.User) UserResponse {
	return UserResponse{
		ID:            u.ID.String(),
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		AccountStatus: u.AccountStatus,
	}
}

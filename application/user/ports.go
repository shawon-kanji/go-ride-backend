package user

// PasswordHasher abstracts password hashing and verification.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) error
}

// TokenManager abstracts token generation for auth flows.
type TokenManager interface {
	Generate(userID string, email string) (string, error)
}

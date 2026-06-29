package driver

type SignupRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type DriverResponse struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	AccountStatus   string `json:"account_status"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

type SignupResponse struct {
	Driver DriverResponse `json:"driver"`
}

type LoginResponse struct {
	AccessToken string         `json:"access_token"`
	Driver      DriverResponse `json:"driver"`
}

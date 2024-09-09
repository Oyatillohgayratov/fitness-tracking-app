package models

type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	ID       int            `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Profile  map[string]any `json:"profile,omitempty"`
}

type PasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordResetSubmitRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type UpdateUserRequest struct {
	ID       int            `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Profile  map[string]any `json:"profile,omitempty"`
}

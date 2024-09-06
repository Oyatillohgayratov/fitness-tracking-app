package router

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

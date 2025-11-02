package request

type CommonRequest struct {
	Name string `json:"name" binding:"required"`
}

type CommonListRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UserRequest represents user-related authentication request
type UserRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Username     string `json:"username,omitempty"`
}

package response

import (
	"time"

	"github.com/ryo-arima/circulator/pkg/entity/model"
)

type CommonResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Commons []Common `json:"commons,omitempty"`
}

type CommonListResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Commons []Common `json:"commons"`
	Total   int      `json:"total"`
}

type Common struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	Code      string           `json:"code"`
	Message   string           `json:"message"`
	Token     string           `json:"token,omitempty"`      // Direct token field for simple auth
	TokenPair *model.TokenPair `json:"token_pair,omitempty"` // Token pair for complex auth
	User      *AuthUser        `json:"user,omitempty"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	Code      string           `json:"code"`
	Message   string           `json:"message"`
	TokenPair *model.TokenPair `json:"token_pair,omitempty"`
}

// ValidateResponse represents token validation response
type ValidateResponse struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Valid   bool      `json:"valid"`
	User    *AuthUser `json:"user,omitempty"`
}

// AuthUser represents authenticated user information in response (avoids name clash with entity User)
type AuthUser struct {
	ID       uint   `json:"id"`
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// UserResponse represents user-related operations response
type UserResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    *User  `json:"data,omitempty"`
}

type UserListResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    []User `json:"data,omitempty"`
	Total   int    `json:"total"`
}

// User represents user information in response
type User struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// RegisterAgentResponse represents agent registration response
type RegisterAgentResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	AgentID string `json:"agent_id"`
	Status  string `json:"status"`
}

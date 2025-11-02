package model

import "time"

type Common struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	UUID      string     `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	Name      string     `gorm:"type:varchar(255)" json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (Common) TableName() string {
	return "commons"
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	Jti       string `json:"jti"`
	UserID    uint   `json:"user_id"`
	UUID      string `json:"uuid"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role,omitempty"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

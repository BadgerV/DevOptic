package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash *string   `json:"-" db:"password_hash"` // Hidden from JSON, nullable for SSO
	Provider     string    `json:"provider" db:"provider"`
	ProviderID   *string   `json:"provider_id,omitempty" db:"provider_id"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Session struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash  string    `json:"-" db:"token_hash"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	LastUsedAt time.Time `json:"last_used_at" db:"last_used_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}

// type AuthProvider interface {
//     Authenticate(credentials interface{}) (*User, error)
//     ValidateToken(token string) (*User, error)
//     GetProviderName() string
// }

type AuthContext struct {
	User        *User
	Permissions []string
	IsAdmin     bool
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByProviderID(ctx context.Context, provider, providerID string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int) ([]User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error
    SetDeliveryEmail(ctx context.Context, userID uuid.UUID, deliveryEmail string) error
    GetDeliveryEmail(ctx context.Context, userID uuid.UUID) (string, error)
    
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error)
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]Session, error)
	UpdateSessionLastUsed(ctx context.Context, id uuid.UUID, lastUsed time.Time) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error
	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

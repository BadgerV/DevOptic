package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Defining Error Standards
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrProviderNotFound   = errors.New("authentication provider not found")
)

type AuthProvider interface {
	Authenticate(ctx context.Context, credentials *LoginRequest) (*User, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
	GetProviderName() string
	SupportsRegistration() bool
	Register(ctx context.Context, userData *RegisterRequest) (*User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	SetDeliveryEmail(ctx context.Context, userID uuid.UUID, deliveryEmail string) error
	GetDeliveryEmail(ctx context.Context, userID uuid.UUID) (string, error)
}

type AuthProviderManager struct {
	providers map[string]AuthProvider
}

func NewAuthProviderManager() *AuthProviderManager {
	return &AuthProviderManager{
		providers: make(map[string]AuthProvider),
	}
}

func (apm *AuthProviderManager) RegisterProvider(provider AuthProvider) {
	apm.providers[provider.GetProviderName()] = provider
}

func (apm *AuthProviderManager) GetProvider(name string) (AuthProvider, error) {
	provider, exists := apm.providers[name]

	if !exists {
		return nil, ErrProviderNotFound
	}

	return provider, nil
}

func (apm *AuthProviderManager) GetAllProviders() map[string]AuthProvider {
	return apm.providers
}

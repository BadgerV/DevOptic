// package auth

// import (
// 	"context"
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"time"

// 	"github.com/google/uuid"
// )

// type Service struct {
// 	userRepo        UserRepository
// 	sessionRepo     SessionRepository
// 	providerManager *AuthProviderManager
// 	tokenExpiry     time.Duration
// }

// func NewService(userRepo UserRepository, sessionRepo SessionRepository) *Service {
// 	service := &Service{
// 		userRepo:        userRepo,
// 		sessionRepo:     sessionRepo,
// 		providerManager: NewAuthProviderManager(),
// 		tokenExpiry:     24 * time.Hour, // Default 24 hours
// 	}

// 	// Register local provider
// 	localProvider := NewLocalProvider(userRepo)
// 	service.providerManager.RegisterProvider(localProvider)

// 	return service
// }

// func (s *Service) Login(ctx context.Context, provider string, credentials interface{}) (*LoginResponse, error) {
// 	authProvider, err := s.providerManager.GetProvider(provider)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user, err := authProvider.Authenticate(ctx, credentials)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create session
// 	token, err := s.generateToken()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate token: %w", err)
// 	}

// 	session := &Session{
// 		ID:         uuid.New(),
// 		UserID:     user.ID,
// 		TokenHash:  s.hashToken(token),
// 		ExpiresAt:  time.Now().Add(s.tokenExpiry),
// 		CreatedAt:  time.Now(),
// 		LastUsedAt: time.Now(),
// 	}

// 	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
// 		return nil, fmt.Errorf("failed to create session: %w", err)
// 	}

// 	return &LoginResponse{
// 		Token:     token,
// 		ExpiresAt: session.ExpiresAt,
// 		User:      *user,
// 	}, nil
// }

// func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
// 	return s.userRepo.GetUserByID(ctx, id)
// }

// func (s *Service) ValidateToken(ctx context.Context, token string) (*AuthContext, error) {
// 	tokenHash := s.hashToken(token)

// 	session, err := s.sessionRepo.GetSessionByTokenHash(ctx, tokenHash)
// 	if err != nil {
// 		return nil, ErrInvalidCredentials
// 	}

// 	if time.Now().After(session.ExpiresAt) {
// 		s.sessionRepo.DeleteSession(ctx, session.ID)
// 		return nil, ErrInvalidCredentials
// 	}

// 	user, err := s.userRepo.GetUserByID(ctx, session.UserID)
// 	if err != nil {
// 		return nil, ErrUserNotFound
// 	}

// 	if !user.IsActive {
// 		return nil, ErrInvalidCredentials
// 	}

// 	// Update last used time
// 	s.sessionRepo.UpdateSessionLastUsed(ctx, session.ID, time.Now())

// 	return &AuthContext{
// 		User: user,
// 	}, nil
// }

// func (s *Service) Logout(ctx context.Context, token string) error {
// 	tokenHash := s.hashToken(token)
// 	return s.sessionRepo.DeleteSessionByTokenHash(ctx, tokenHash)
// }

// func (s *Service) RegisterProvider(provider AuthProvider) {
// 	s.providerManager.RegisterProvider(provider)
// }

// func (s *Service) generateToken() (string, error) {
// 	bytes := make([]byte, 32)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(bytes), nil
// }

// func (s *Service) hashToken(token string) string {
// 	// Simple hash for demo - use proper hashing in production
// 	return fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
// }

package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"crypto/hmac"

	"github.com/google/uuid"
)

type Service struct {
	userRepo        UserRepository
	sessionRepo     SessionRepository
	providerManager *AuthProviderManager
	tokenExpiry     time.Duration
}

// func NewService(userRepo UserRepository) *Service {
// 	service := &Service{
// 		userRepo:        userRepo,
// 		providerManager: NewAuthProviderManager(),
// 		tokenExpiry:     24 * time.Hour, // Default 24 hours
// 	}

// 	// Register local provider
// 	localProvider := NewLocalProvider(userRepo)
// 	service.providerManager.RegisterProvider(localProvider)

// 	return service
// }

func NewService(userRepo UserRepository, sessionRepo SessionRepository) *Service {
	service := &Service{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		providerManager: NewAuthProviderManager(),
		tokenExpiry:     24 * time.Hour, // Default 24 hours
	}

	// Register local provider
	localProvider := NewLocalProvider(userRepo)
	service.providerManager.RegisterProvider(localProvider)

	return service
}

func (s *Service) Login(ctx context.Context, provider string, credentials *LoginRequest) (*LoginResponse, error) {
	authProvider, err := s.providerManager.GetProvider(provider)
	if err != nil {
		return nil, err
	}

	user, err := authProvider.Authenticate(ctx, credentials)
	if err != nil {
		return nil, err
	}

	// Create session
	token, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session := &Session{
		ID:         uuid.New(),
		UserID:     user.ID,
		TokenHash:  s.hashToken(token),
		ExpiresAt:  time.Now().Add(s.tokenExpiry),
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginResponse{
		Token:     token,
		ExpiresAt: session.ExpiresAt,
		User:      *user,
	}, nil
}

func (s *Service) Register(ctx context.Context, provider string, userData *RegisterRequest) (*LoginResponse, error) {
	// 1. Get the provider (e.g. LocalProvider)
	authProvider, err := s.providerManager.GetProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// 2. Call the provider-specific Register
	user, err := authProvider.Register(ctx, userData)
	if err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	// 3. Generate a session token for the new user
	token, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// (Optional) Persist session in DB
	session := &Session{
		ID:         uuid.New(),
		UserID:     user.ID,
		TokenHash:  s.hashToken(token),
		ExpiresAt:  time.Now().Add(s.tokenExpiry),
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}
	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	fmt.Println("hasshdetoken is - ", session.TokenHash)

	// 4. Return login response (same shape as Login)
	return &LoginResponse{
		Token:     token,
		ExpiresAt: session.ExpiresAt,
		User:      *user,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *Service) ValidateToken(ctx context.Context, token string) (*AuthContext, error) {
	fmt.Println("Token entering validate token function", token)
	tokenHash := s.hashToken(token)

	session, err := s.sessionRepo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if time.Now().After(session.ExpiresAt) {
		s.sessionRepo.DeleteSession(ctx, session.ID)
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetUserByID(ctx, session.UserID)

	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Update last used time
	s.sessionRepo.UpdateSessionLastUsed(ctx, session.ID, time.Now())

	return &AuthContext{
		User: user,
	}, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	tokenHash := s.hashToken(token)
	return s.sessionRepo.DeleteSessionByTokenHash(ctx, tokenHash)
}

func (s *Service) RegisterProvider(provider AuthProvider) {
	s.providerManager.RegisterProvider(provider)
}

func (s *Service) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// func (s *Service) hashToken(token string) string {
// 	// Simple hash for demo - use proper hashing in production
// 	return fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
// }

func (s *Service) hashToken(token string) string {
	secret := []byte("super-strong-secret-key") // store securely (env var, Vault, etc.)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

// ChangePassword updates a user's password via the specified provider.
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
    authProvider, err := s.providerManager.GetProvider("local")
    if err != nil {
        return err
    }

    err = authProvider.ChangePassword(ctx, userID, oldPassword, newPassword)
    if err != nil {
        return err
    }

    // Optionally invalidate existing sessions for security
    // err = s.sessionRepo.DeleteSessionsByUserID(ctx, userID)
    // if err != nil {
    //     return fmt.Errorf("failed to invalidate sessions: %w", err)
    // }

    return nil
}

// SetDeliveryEmail sets the delivery email for a user via the specified provider.
func (s *Service) SetDeliveryEmail(ctx context.Context, userID uuid.UUID, deliveryEmail string) error {
    authProvider, err := s.providerManager.GetProvider("local")
    if err != nil {
        return err
    }

    err = authProvider.SetDeliveryEmail(ctx, userID, deliveryEmail)
    if err != nil {
        return err
    }

    return nil
}

// GetDeliveryEmail retrieves the delivery email for a user via the specified provider.
func (s *Service) GetDeliveryEmail(ctx context.Context, userID uuid.UUID) (string, error) {
    authProvider, err := s.providerManager.GetProvider("local")
    if err != nil {
        return "", err
    }

    deliveryEmail, err := authProvider.GetDeliveryEmail(ctx, userID)
    if err != nil {
        return "", err
    }

    return deliveryEmail, nil
}
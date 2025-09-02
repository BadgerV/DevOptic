package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type LocalProvider struct {
	userRepo UserRepository
}

func NewLocalProvider(userRepo UserRepository) *LocalProvider {
	return &LocalProvider{
		userRepo: userRepo,
	}
}

func (lp *LocalProvider) GetProviderName() string {
	return "local"
}

func (lp *LocalProvider) SupportsRegistration() bool {
	return true
}

func (lp *LocalProvider) Authenticate(ctx context.Context, creds *LoginRequest) (*User, error) {
	log.Println(creds)

	if creds == nil {
		return nil, ErrInvalidCredentials
	}

	user, err := lp.userRepo.GetUserByUsername(ctx, creds.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	if !lp.verifyPassword(creds.Password, *user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	return user, nil
}

func (lp *LocalProvider) Register(ctx context.Context, regData *RegisterRequest) (*User, error) {
	log.Println(regData)

	if regData == nil {
		return nil, errors.New("registration data cannot be nil")
	}

	hashedPassword, err := lp.hashPassword(regData.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		ID:           uuid.New(),
		Username:     regData.Username,
		Email:        regData.Email,
		PasswordHash: &hashedPassword,
		Provider:     "local",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return lp.userRepo.CreateUser(ctx, user)
}

func (lp *LocalProvider) ValidateToken(ctx context.Context, token string) (*User, error) {
	// This would typically validate JWT tokens or session tokens
	// Implementation depends on your token strategy
	return nil, fmt.Errorf("not implemented")
}

func (lp *LocalProvider) verifyPassword(password, hashedPassword string) bool {
	parts := strings.Split(hashedPassword, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return subtle.ConstantTimeCompare(hash, comparisonHash) == 1
}

func (lp *LocalProvider) hashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return fmt.Sprintf("%s:%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash)), nil
}

// ChangePassword updates a user's password after verifying the old password.
func (lp *LocalProvider) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
    log.Println("Changing password for user ID:", userID)
    
    user, err := lp.userRepo.GetUserByID(ctx, userID)
    if err != nil {
        return ErrUserNotFound
    }
    
    if user.PasswordHash == nil {
        return ErrInvalidCredentials
    }
    
    // Verify old password
    if !lp.verifyPassword(oldPassword, *user.PasswordHash) {
        return ErrInvalidCredentials
    }
    
    // Hash new password
    newPasswordHash, err := lp.hashPassword(newPassword)
    if err != nil {
        return fmt.Errorf("failed to hash new password: %w", err)
    }
    
    // Update password in the repository
    err = lp.userRepo.UpdatePassword(ctx, userID, newPasswordHash)
    if err != nil {
        return err
    }
    
    return nil
}

// SetDeliveryEmail sets the delivery email for a user.
func (lp *LocalProvider) SetDeliveryEmail(ctx context.Context, userID uuid.UUID, deliveryEmail string) error {
    log.Println("Setting delivery email for user ID:", userID)
    
    user, err := lp.userRepo.GetUserByID(ctx, userID)
    if err != nil {
        return ErrUserNotFound
    }
    
    if !user.IsActive {
        return errors.New("user account is not active")
    }
    
    err = lp.userRepo.SetDeliveryEmail(ctx, userID, deliveryEmail)
    if err != nil {
        return err
    }
    
    return nil
}

// GetDeliveryEmail retrieves the delivery email for a user, falling back to their primary email if unset.
func (lp *LocalProvider) GetDeliveryEmail(ctx context.Context, userID uuid.UUID) (string, error) {
    log.Println("Fetching delivery email for user ID:", userID)
    
    user, err := lp.userRepo.GetUserByID(ctx, userID)
    if err != nil {
        return "", ErrUserNotFound
    }
    
    if !user.IsActive {
        return "", errors.New("user account is not active")
    }
    
    deliveryEmail, err := lp.userRepo.GetDeliveryEmail(ctx, userID)
    if err != nil {
        return "", err
    }
    
    return deliveryEmail, nil
}
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/badgerv/monitoring-api/internal/storage" // Update this import path to match your project
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresUserRepository struct {
	db *storage.DB
}

func NewPostgresUserRepository(db *storage.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	query := `
        INSERT INTO users (id, username, email, password_hash, provider, provider_id, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, username, email, password_hash, provider, provider_id, is_active, created_at, updated_at`

	var created User
	err := r.db.Pool.QueryRow(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.Provider, user.ProviderID, user.IsActive,
		user.CreatedAt, user.UpdatedAt).
		Scan(&created.ID, &created.Username, &created.Email, &created.PasswordHash,
			&created.Provider, &created.ProviderID, &created.IsActive,
			&created.CreatedAt, &created.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &created, nil
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, provider, provider_id, is_active, created_at, updated_at 
              FROM users WHERE username = $1 AND is_active = true`

	var user User
	err := r.db.Pool.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Provider, &user.ProviderID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `SELECT id, username, email, password_hash, provider, provider_id, is_active, created_at, updated_at 
              FROM users WHERE id = $1`

	var user User
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Provider, &user.ProviderID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password_hash, provider, provider_id, is_active, created_at, updated_at 
              FROM users WHERE email = $1 AND is_active = true`

	var user User
	err := r.db.Pool.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Provider, &user.ProviderID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserByProviderID(ctx context.Context, provider, providerID string) (*User, error) {
	query := `SELECT id, username, email, password_hash, provider, provider_id, is_active, created_at, updated_at 
              FROM users WHERE provider = $1 AND provider_id = $2 AND is_active = true`

	var user User
	err := r.db.Pool.QueryRow(ctx, query, provider, providerID).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Provider, &user.ProviderID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *User) error {
	query := `
        UPDATE users 
        SET username = $2, email = $3, password_hash = $4, provider = $5, 
            provider_id = $6, is_active = $7, updated_at = $8
        WHERE id = $1`

	user.UpdatedAt = time.Now()

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.Provider, user.ProviderID, user.IsActive, user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = $2 WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) ListUsers(ctx context.Context, limit, offset int) ([]User, error) {
	query := `SELECT id, username, email, password_hash, provider, provider_id, is_active, created_at, updated_at 
              FROM users WHERE is_active = true ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Provider, &user.ProviderID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating users: %w", rows.Err())
	}

	return users, nil
}

type PostgresSessionRepository struct {
	db *storage.DB
}

func NewPostgresSessionRepository(db *storage.DB) SessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) CreateSession(ctx context.Context, session *Session) error {
	query := `
        INSERT INTO user_sessions (id, user_id, token_hash, expires_at, created_at, last_used_at)
        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Pool.Exec(ctx, query,
		session.ID, session.UserID, session.TokenHash,
		session.ExpiresAt, session.CreatedAt, session.LastUsedAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *PostgresSessionRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*Session, error) {
	fmt.Println("Token to be validated", tokenHash)
	query := `SELECT id, user_id, token_hash, expires_at, created_at, last_used_at 
              FROM user_sessions WHERE token_hash = $1 AND expires_at > NOW() AT TIME ZONE 'UTC'`

	var session Session
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).
		Scan(&session.ID, &session.UserID, &session.TokenHash,
			&session.ExpiresAt, &session.CreatedAt, &session.LastUsedAt)
	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	fmt.Println(session)

	return &session, nil
}

func (r *PostgresSessionRepository) UpdateSessionLastUsed(ctx context.Context, id uuid.UUID, lastUsed time.Time) error {
	query := `UPDATE user_sessions SET last_used_at = $2 WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query, id, lastUsed)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

func (r *PostgresSessionRepository) DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM user_sessions WHERE token_hash = $1`

	_, err := r.db.Pool.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r *PostgresSessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	query := `SELECT id, user_id, token_hash, expires_at, created_at, last_used_at 
              FROM user_sessions WHERE id = $1`

	var session Session
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&session.ID, &session.UserID, &session.TokenHash,
			&session.ExpiresAt, &session.CreatedAt, &session.LastUsedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (r *PostgresSessionRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	query := `SELECT id, user_id, token_hash, expires_at, created_at, last_used_at 
              FROM user_sessions WHERE user_id = $1 AND expires_at > NOW() ORDER BY last_used_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var session Session
		err := rows.Scan(&session.ID, &session.UserID, &session.TokenHash,
			&session.ExpiresAt, &session.CreatedAt, &session.LastUsedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating user sessions: %w", rows.Err())
	}

	return sessions, nil
}

func (r *PostgresSessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r *PostgresSessionRepository) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE user_id = $1`

	_, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}

func (r *PostgresSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE expires_at <= NOW()`

	_, err := r.db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}

// Repository function to update password (unchanged, as it was correct)
func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error {
	fmt.Println("Updating password for user ID:", userID)
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`

	_, err := r.db.Pool.Exec(ctx, query, newPasswordHash, userID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Repository function to set delivery_email
func (r *PostgresUserRepository) SetDeliveryEmail(ctx context.Context, userID uuid.UUID, deliveryEmail string) error {
	fmt.Println("Setting delivery email for user ID:", userID)
	query := `UPDATE users SET delivery_email = $1 WHERE id = $2`

	_, err := r.db.Pool.Exec(ctx, query, deliveryEmail, userID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to set delivery email: %w", err)
	}

	return nil
}

// Repository function to fetch delivery_email (fallback to email if null)
func (r *PostgresUserRepository) GetDeliveryEmail(ctx context.Context, userID uuid.UUID) (string, error) {
	fmt.Println("Fetching delivery email for user ID:", userID)
	query := `SELECT COALESCE(delivery_email, email) AS effective_email FROM users WHERE id = $1`

	var effectiveEmail string
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&effectiveEmail)
	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get delivery email: %w", err)
	}

	fmt.Println("\n\n\n\nDelivery Email:", effectiveEmail)

	return effectiveEmail, nil
}

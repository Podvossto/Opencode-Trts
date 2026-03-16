// Purpose: Auth domain — repository layer (database queries)
// DB tables: users, otp_logs
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRow represents a row from the users table relevant to auth
type UserRow struct {
	ID                 uuid.UUID
	Email              string
	PasswordHash       string
	FirstName          string
	LastName           string
	Role               string
	IsActive           bool
	MustChangePassword bool
}

// OTPRow represents a row from the otp_logs table
type OTPRow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	OTPCode   string
	IsUsed    bool
	ExpiresAt time.Time
}

// RepositoryInterface defines the contract for auth data access (enables mocking in tests)
type RepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*UserRow, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserRow, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string, clearMustChange bool) error
	CreateOTP(ctx context.Context, userID uuid.UUID, otpCode string, expiresAt time.Time) error
	GetLatestValidOTP(ctx context.Context, userID uuid.UUID) (*OTPRow, error)
	MarkOTPUsed(ctx context.Context, otpID uuid.UUID) error
}

// Repository handles auth-related database operations
type Repository struct {
	db *pgxpool.Pool
}

// compile-time check: Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

// NewRepository creates a new auth repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetUserByEmail retrieves a user by email for login
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*UserRow, error) {
	var u UserRow
	err := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, first_name, last_name, role, is_active, must_change_password
		 FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.IsActive, &u.MustChangePassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*UserRow, error) {
	var u UserRow
	err := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, first_name, last_name, role, is_active, must_change_password
		 FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.IsActive, &u.MustChangePassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

// UpdatePassword updates a user's password hash and optionally clears the must_change_password flag
func (r *Repository) UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string, clearMustChange bool) error {
	query := `UPDATE users SET password_hash = $1, updated_at = now() WHERE id = $2`
	if clearMustChange {
		query = `UPDATE users SET password_hash = $1, must_change_password = false, updated_at = now() WHERE id = $2`
	}
	_, err := r.db.Exec(ctx, query, newHash, userID)
	return err
}

// CreateOTP inserts a new OTP log entry, invalidating any existing unused OTPs for the user
func (r *Repository) CreateOTP(ctx context.Context, userID uuid.UUID, otpCode string, expiresAt time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Invalidate existing unused OTPs
	_, err = tx.Exec(ctx,
		`UPDATE otp_logs SET is_used = true WHERE user_id = $1 AND is_used = false`, userID)
	if err != nil {
		return fmt.Errorf("invalidate old OTPs: %w", err)
	}

	// Insert new OTP
	_, err = tx.Exec(ctx,
		`INSERT INTO otp_logs (user_id, otp_code, expires_at) VALUES ($1, $2, $3)`,
		userID, otpCode, expiresAt)
	if err != nil {
		return fmt.Errorf("insert OTP: %w", err)
	}

	return tx.Commit(ctx)
}

// GetLatestValidOTP retrieves the latest unused, unexpired OTP for a user
func (r *Repository) GetLatestValidOTP(ctx context.Context, userID uuid.UUID) (*OTPRow, error) {
	var o OTPRow
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, otp_code, is_used, expires_at
		 FROM otp_logs
		 WHERE user_id = $1 AND is_used = false AND expires_at > now()
		 ORDER BY created_at DESC LIMIT 1`, userID).
		Scan(&o.ID, &o.UserID, &o.OTPCode, &o.IsUsed, &o.ExpiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get latest valid OTP: %w", err)
	}
	return &o, nil
}

// MarkOTPUsed marks an OTP record as used
func (r *Repository) MarkOTPUsed(ctx context.Context, otpID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE otp_logs SET is_used = true WHERE id = $1`, otpID)
	return err
}

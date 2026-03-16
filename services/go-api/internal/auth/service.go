// Purpose: Auth domain — service layer (business logic)
// Ref: API Contract — POST /api/v1/auth/login, POST /api/v1/auth/logout, OTP flow
// DB tables: users, otp_logs
package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/tigersoft/ats-go-api/config"
	pkgjwt "github.com/tigersoft/ats-go-api/pkg/jwt"
	pkgredis "github.com/tigersoft/ats-go-api/pkg/redis"
)

// LoginRequest is the parsed body for POST /api/v1/auth/login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse returned on successful login
type LoginResponse struct {
	AccessToken        string `json:"access_token"`
	ExpiresIn          int    `json:"expires_in"` // seconds
	Role               string `json:"role"`
	MustChangePassword bool   `json:"must_change_password"`
}

// ChangePasswordRequest is the body for PUT /api/v1/auth/password/change
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// OTPRequestBody for POST /api/v1/auth/otp/request
type OTPRequestBody struct {
	Email string `json:"email" binding:"required,email"`
}

// OTPVerifyBody for POST /api/v1/auth/otp/verify
type OTPVerifyBody struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// PasswordResetBody for POST /api/v1/auth/password/reset
type PasswordResetBody struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// RedisBlacklister abstracts the Redis blacklist operations (enables mocking in tests)
type RedisBlacklister interface {
	BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error
}

// redisBlacklisterAdapter adapts *goredis.Client to the RedisBlacklister interface
type redisBlacklisterAdapter struct {
	rdb *goredis.Client
}

func (a *redisBlacklisterAdapter) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	return pkgredis.BlacklistToken(ctx, a.rdb, jti, ttl)
}

// Service implements the auth business logic
type Service struct {
	repo        RepositoryInterface
	blacklister RedisBlacklister
	cfg         *config.Config
}

// NewService creates a new auth service
func NewService(repo *Repository, rdb *goredis.Client, cfg *config.Config) *Service {
	return &Service{repo: repo, blacklister: &redisBlacklisterAdapter{rdb: rdb}, cfg: cfg}
}

// NewServiceWithDeps creates a new auth service with explicit interface deps (for testing)
func NewServiceWithDeps(repo RepositoryInterface, blacklister RedisBlacklister, cfg *config.Config) *Service {
	return &Service{repo: repo, blacklister: blacklister, cfg: cfg}
}

// Login validates credentials and issues a JWT
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT
	token, _, err := pkgjwt.Generate(user.ID.String(), user.Role, s.cfg.JWTSecret, s.cfg.JWTExpiry)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &LoginResponse{
		AccessToken:        token,
		ExpiresIn:          s.cfg.JWTExpiry * 60,
		Role:               user.Role,
		MustChangePassword: user.MustChangePassword,
	}, nil
}

// Logout blacklists the current JWT
func (s *Service) Logout(ctx context.Context, jti string, ttl time.Duration) error {
	return s.blacklister.BlacklistToken(ctx, jti, ttl)
}

// ChangePassword validates old password and updates to new password
func (s *Service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return fmt.Errorf("change password: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Prevent setting same password
	if req.OldPassword == req.NewPassword {
		return ErrSamePassword
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.repo.UpdatePassword(ctx, uid, string(newHash), true)
}

// RequestOTP generates a 6-digit OTP and stores it (email sending is handled by handler/caller)
func (s *Service) RequestOTP(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("request OTP: %w", err)
	}
	if user == nil {
		// Return success silently for security (no email enumeration)
		return "", nil
	}

	// Admin cannot use OTP reset (SRS FRA07)
	if user.Role == "admin" {
		return "", ErrAdminNoOTP
	}

	// Generate 6-digit OTP
	otp, err := generateOTP(6)
	if err != nil {
		return "", fmt.Errorf("generate OTP: %w", err)
	}

	// Store OTP with 5-minute expiry
	expiresAt := time.Now().Add(5 * time.Minute)
	if err := s.repo.CreateOTP(ctx, user.ID, otp, expiresAt); err != nil {
		return "", fmt.Errorf("store OTP: %w", err)
	}

	return otp, nil
}

// VerifyOTPAndResetPassword verifies OTP and resets the password in one step
func (s *Service) VerifyOTPAndResetPassword(ctx context.Context, email, otp, newPassword string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("verify OTP: %w", err)
	}
	if user == nil {
		return ErrInvalidOTP
	}

	// Get latest valid OTP
	otpRow, err := s.repo.GetLatestValidOTP(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("get OTP: %w", err)
	}
	if otpRow == nil || otpRow.OTPCode != otp {
		return ErrInvalidOTP
	}

	// Mark OTP as used
	if err := s.repo.MarkOTPUsed(ctx, otpRow.ID); err != nil {
		return fmt.Errorf("mark OTP used: %w", err)
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.repo.UpdatePassword(ctx, user.ID, string(newHash), false)
}

// generateOTP creates a cryptographically secure n-digit numeric OTP
func generateOTP(length int) (string, error) {
	digits := make([]byte, length)
	for i := range digits {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		digits[i] = byte('0') + byte(n.Int64())
	}
	return string(digits), nil
}

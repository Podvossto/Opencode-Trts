package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/tigersoft/ats-go-api/config"
)

// --- Mock Repository ---

type mockRepo struct {
	getUserByEmailFn    func(ctx context.Context, email string) (*UserRow, error)
	getUserByIDFn       func(ctx context.Context, id uuid.UUID) (*UserRow, error)
	updatePasswordFn    func(ctx context.Context, userID uuid.UUID, newHash string, clearMustChange bool) error
	createOTPFn         func(ctx context.Context, userID uuid.UUID, otpCode string, expiresAt time.Time) error
	getLatestValidOTPFn func(ctx context.Context, userID uuid.UUID) (*OTPRow, error)
	markOTPUsedFn       func(ctx context.Context, otpID uuid.UUID) error
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (*UserRow, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *mockRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*UserRow, error) {
	if m.getUserByIDFn != nil {
		return m.getUserByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string, clearMustChange bool) error {
	if m.updatePasswordFn != nil {
		return m.updatePasswordFn(ctx, userID, newHash, clearMustChange)
	}
	return nil
}

func (m *mockRepo) CreateOTP(ctx context.Context, userID uuid.UUID, otpCode string, expiresAt time.Time) error {
	if m.createOTPFn != nil {
		return m.createOTPFn(ctx, userID, otpCode, expiresAt)
	}
	return nil
}

func (m *mockRepo) GetLatestValidOTP(ctx context.Context, userID uuid.UUID) (*OTPRow, error) {
	if m.getLatestValidOTPFn != nil {
		return m.getLatestValidOTPFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockRepo) MarkOTPUsed(ctx context.Context, otpID uuid.UUID) error {
	if m.markOTPUsedFn != nil {
		return m.markOTPUsedFn(ctx, otpID)
	}
	return nil
}

// --- Mock Blacklister ---

type mockBlacklister struct {
	blacklistTokenFn func(ctx context.Context, jti string, ttl time.Duration) error
}

func (m *mockBlacklister) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	if m.blacklistTokenFn != nil {
		return m.blacklistTokenFn(ctx, jti, ttl)
	}
	return nil
}

// --- Test Helpers ---

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret: "test-secret-key-minimum-32-chars!!",
		JWTExpiry: 30,
	}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost) // MinCost for speed in tests
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hash)
}

func activeUser(t *testing.T, email, password, role string) *UserRow {
	t.Helper()
	return &UserRow{
		ID:                 uuid.New(),
		Email:              email,
		PasswordHash:       hashPassword(t, password),
		FirstName:          "Test",
		LastName:           "User",
		Role:               role,
		IsActive:           true,
		MustChangePassword: false,
	}
}

// ==================== Login Tests ====================

func TestLogin_Success(t *testing.T) {
	user := activeUser(t, "hr@test.com", "password123", "hr")
	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			if email == "hr@test.com" {
				return user, nil
			}
			return nil, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	resp, err := svc.Login(context.Background(), LoginRequest{
		Email:    "hr@test.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("Login() returned empty AccessToken")
	}
	if resp.Role != "hr" {
		t.Errorf("Login() Role = %q, want %q", resp.Role, "hr")
	}
	if resp.ExpiresIn != 30*60 {
		t.Errorf("Login() ExpiresIn = %d, want %d", resp.ExpiresIn, 30*60)
	}
}

func TestLogin_InvalidEmail(t *testing.T) {
	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return nil, nil // user not found
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	_, err := svc.Login(context.Background(), LoginRequest{
		Email:    "notfound@test.com",
		Password: "password123",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	user := activeUser(t, "hr@test.com", "correctpassword", "hr")
	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return user, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	_, err := svc.Login(context.Background(), LoginRequest{
		Email:    "hr@test.com",
		Password: "wrongpassword",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestLogin_DisabledAccount(t *testing.T) {
	user := activeUser(t, "hr@test.com", "password123", "hr")
	user.IsActive = false

	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return user, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	_, err := svc.Login(context.Background(), LoginRequest{
		Email:    "hr@test.com",
		Password: "password123",
	})

	if !errors.Is(err, ErrAccountDisabled) {
		t.Errorf("Login() error = %v, want ErrAccountDisabled", err)
	}
}

func TestLogin_MustChangePassword(t *testing.T) {
	user := activeUser(t, "hr@test.com", "password123", "hr")
	user.MustChangePassword = true

	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return user, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	resp, err := svc.Login(context.Background(), LoginRequest{
		Email:    "hr@test.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if !resp.MustChangePassword {
		t.Error("Login() MustChangePassword should be true")
	}
}

func TestLogin_RepoError(t *testing.T) {
	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return nil, errors.New("db connection lost")
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	_, err := svc.Login(context.Background(), LoginRequest{
		Email:    "hr@test.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("Login() should return error on repo failure")
	}
}

// ==================== Logout Tests ====================

func TestLogout_Success(t *testing.T) {
	called := false
	bl := &mockBlacklister{
		blacklistTokenFn: func(ctx context.Context, jti string, ttl time.Duration) error {
			called = true
			if jti != "test-jti-123" {
				t.Errorf("Logout() jti = %q, want %q", jti, "test-jti-123")
			}
			return nil
		},
	}

	svc := NewServiceWithDeps(&mockRepo{}, bl, testConfig())
	err := svc.Logout(context.Background(), "test-jti-123", 30*time.Minute)

	if err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if !called {
		t.Error("Logout() did not call BlacklistToken")
	}
}

func TestLogout_RedisError(t *testing.T) {
	bl := &mockBlacklister{
		blacklistTokenFn: func(ctx context.Context, jti string, ttl time.Duration) error {
			return errors.New("redis down")
		},
	}

	svc := NewServiceWithDeps(&mockRepo{}, bl, testConfig())
	err := svc.Logout(context.Background(), "test-jti", 30*time.Minute)

	if err == nil {
		t.Fatal("Logout() should return error on redis failure")
	}
}

// ==================== ChangePassword Tests ====================

func TestChangePassword_Success(t *testing.T) {
	userID := uuid.New()
	user := &UserRow{
		ID:           userID,
		PasswordHash: hashPassword(t, "oldpassword1"),
		IsActive:     true,
	}

	var updatedHash string
	var clearedMustChange bool

	repo := &mockRepo{
		getUserByIDFn: func(ctx context.Context, id uuid.UUID) (*UserRow, error) {
			if id == userID {
				return user, nil
			}
			return nil, nil
		},
		updatePasswordFn: func(ctx context.Context, uid uuid.UUID, newHash string, clearMC bool) error {
			updatedHash = newHash
			clearedMustChange = clearMC
			return nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	err := svc.ChangePassword(context.Background(), userID.String(), ChangePasswordRequest{
		OldPassword: "oldpassword1",
		NewPassword: "newpassword2",
	})

	if err != nil {
		t.Fatalf("ChangePassword() error = %v", err)
	}
	if updatedHash == "" {
		t.Error("ChangePassword() did not call UpdatePassword")
	}
	if !clearedMustChange {
		t.Error("ChangePassword() should clear must_change_password flag")
	}

	// Verify new hash is valid
	if err := bcrypt.CompareHashAndPassword([]byte(updatedHash), []byte("newpassword2")); err != nil {
		t.Error("ChangePassword() new hash does not match new password")
	}
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	userID := uuid.New()
	user := &UserRow{
		ID:           userID,
		PasswordHash: hashPassword(t, "correctpassword"),
		IsActive:     true,
	}

	repo := &mockRepo{
		getUserByIDFn: func(ctx context.Context, id uuid.UUID) (*UserRow, error) {
			return user, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	err := svc.ChangePassword(context.Background(), userID.String(), ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword2",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("ChangePassword() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestChangePassword_SamePassword(t *testing.T) {
	userID := uuid.New()
	user := &UserRow{
		ID:           userID,
		PasswordHash: hashPassword(t, "samepassword1"),
		IsActive:     true,
	}

	repo := &mockRepo{
		getUserByIDFn: func(ctx context.Context, id uuid.UUID) (*UserRow, error) {
			return user, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	err := svc.ChangePassword(context.Background(), userID.String(), ChangePasswordRequest{
		OldPassword: "samepassword1",
		NewPassword: "samepassword1",
	})

	if !errors.Is(err, ErrSamePassword) {
		t.Errorf("ChangePassword() error = %v, want ErrSamePassword", err)
	}
}

func TestChangePassword_UserNotFound(t *testing.T) {
	repo := &mockRepo{
		getUserByIDFn: func(ctx context.Context, id uuid.UUID) (*UserRow, error) {
			return nil, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	err := svc.ChangePassword(context.Background(), uuid.New().String(), ChangePasswordRequest{
		OldPassword: "oldpassword1",
		NewPassword: "newpassword2",
	})

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("ChangePassword() error = %v, want ErrUserNotFound", err)
	}
}

func TestChangePassword_InvalidUserID(t *testing.T) {
	svc := NewServiceWithDeps(&mockRepo{}, &mockBlacklister{}, testConfig())
	err := svc.ChangePassword(context.Background(), "not-a-uuid", ChangePasswordRequest{
		OldPassword: "oldpassword1",
		NewPassword: "newpassword2",
	})

	if err == nil {
		t.Fatal("ChangePassword() should return error for invalid UUID")
	}
}

// ==================== RequestOTP Tests ====================

func TestRequestOTP_Success(t *testing.T) {
	userID := uuid.New()
	user := &UserRow{
		ID:    userID,
		Email: "hr@test.com",
		Role:  "hr",
	}

	var storedOTP string
	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return user, nil
		},
		createOTPFn: func(ctx context.Context, uid uuid.UUID, otpCode string, expiresAt time.Time) error {
			storedOTP = otpCode
			if uid != userID {
				t.Errorf("CreateOTP called with uid %v, want %v", uid, userID)
			}
			// Verify OTP is 6 digits
			if len(otpCode) != 6 {
				t.Errorf("OTP length = %d, want 6", len(otpCode))
			}
			for _, c := range otpCode {
				if c < '0' || c > '9' {
					t.Errorf("OTP contains non-digit char: %c", c)
				}
			}
			// Verify expiry is ~5 minutes from now
			diff := expiresAt.Sub(time.Now())
			if diff < 4*time.Minute || diff > 6*time.Minute {
				t.Errorf("OTP expiry off: %v", diff)
			}
			return nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	otp, err := svc.RequestOTP(context.Background(), "hr@test.com")

	if err != nil {
		t.Fatalf("RequestOTP() error = %v", err)
	}
	if otp == "" {
		t.Error("RequestOTP() returned empty OTP")
	}
	if otp != storedOTP {
		t.Errorf("RequestOTP() returned OTP %q but stored %q", otp, storedOTP)
	}
}

func TestRequestOTP_AdminBlocked(t *testing.T) {
	user := &UserRow{
		ID:    uuid.New(),
		Email: "admin@test.com",
		Role:  "admin",
	}

	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return user, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	_, err := svc.RequestOTP(context.Background(), "admin@test.com")

	if !errors.Is(err, ErrAdminNoOTP) {
		t.Errorf("RequestOTP() error = %v, want ErrAdminNoOTP", err)
	}
}

func TestRequestOTP_UnknownEmail_SilentSuccess(t *testing.T) {
	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return nil, nil // user not found
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	otp, err := svc.RequestOTP(context.Background(), "unknown@test.com")

	if err != nil {
		t.Fatalf("RequestOTP() should succeed silently for unknown email, got error = %v", err)
	}
	if otp != "" {
		t.Errorf("RequestOTP() should return empty OTP for unknown email, got %q", otp)
	}
}

func TestRequestOTP_RepoError(t *testing.T) {
	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	_, err := svc.RequestOTP(context.Background(), "hr@test.com")

	if err == nil {
		t.Fatal("RequestOTP() should return error on repo failure")
	}
}

// ==================== VerifyOTPAndResetPassword Tests ====================

func TestVerifyOTPAndResetPassword_Success(t *testing.T) {
	userID := uuid.New()
	otpID := uuid.New()
	user := &UserRow{ID: userID, Email: "hr@test.com"}
	otpRow := &OTPRow{ID: otpID, UserID: userID, OTPCode: "123456", IsUsed: false, ExpiresAt: time.Now().Add(5 * time.Minute)}

	var markedUsed bool
	var updatedHash string
	var clearMC bool

	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return user, nil
		},
		getLatestValidOTPFn: func(ctx context.Context, uid uuid.UUID) (*OTPRow, error) {
			return otpRow, nil
		},
		markOTPUsedFn: func(ctx context.Context, oid uuid.UUID) error {
			markedUsed = true
			if oid != otpID {
				t.Errorf("MarkOTPUsed called with %v, want %v", oid, otpID)
			}
			return nil
		},
		updatePasswordFn: func(ctx context.Context, uid uuid.UUID, newHash string, clearMustChange bool) error {
			updatedHash = newHash
			clearMC = clearMustChange
			return nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	err := svc.VerifyOTPAndResetPassword(context.Background(), "hr@test.com", "123456", "newpassword2")

	if err != nil {
		t.Fatalf("VerifyOTPAndResetPassword() error = %v", err)
	}
	if !markedUsed {
		t.Error("VerifyOTPAndResetPassword() should mark OTP as used")
	}
	if updatedHash == "" {
		t.Error("VerifyOTPAndResetPassword() should update password")
	}
	if clearMC {
		t.Error("VerifyOTPAndResetPassword() should NOT clear must_change_password (reset sets it to keep default)")
	}

	// Verify new hash
	if err := bcrypt.CompareHashAndPassword([]byte(updatedHash), []byte("newpassword2")); err != nil {
		t.Error("VerifyOTPAndResetPassword() new hash does not match new password")
	}
}

func TestVerifyOTPAndResetPassword_InvalidOTP(t *testing.T) {
	userID := uuid.New()
	user := &UserRow{ID: userID, Email: "hr@test.com"}
	otpRow := &OTPRow{ID: uuid.New(), UserID: userID, OTPCode: "123456"}

	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return user, nil
		},
		getLatestValidOTPFn: func(ctx context.Context, uid uuid.UUID) (*OTPRow, error) {
			return otpRow, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	err := svc.VerifyOTPAndResetPassword(context.Background(), "hr@test.com", "999999", "newpassword2")

	if !errors.Is(err, ErrInvalidOTP) {
		t.Errorf("VerifyOTPAndResetPassword() error = %v, want ErrInvalidOTP", err)
	}
}

func TestVerifyOTPAndResetPassword_NoValidOTP(t *testing.T) {
	userID := uuid.New()
	user := &UserRow{ID: userID, Email: "hr@test.com"}

	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return user, nil
		},
		getLatestValidOTPFn: func(ctx context.Context, uid uuid.UUID) (*OTPRow, error) {
			return nil, nil // no valid OTP
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	err := svc.VerifyOTPAndResetPassword(context.Background(), "hr@test.com", "123456", "newpassword2")

	if !errors.Is(err, ErrInvalidOTP) {
		t.Errorf("VerifyOTPAndResetPassword() error = %v, want ErrInvalidOTP", err)
	}
}

func TestVerifyOTPAndResetPassword_UserNotFound(t *testing.T) {
	repo := &mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return nil, nil
		},
	}

	svc := NewServiceWithDeps(repo, &mockBlacklister{}, testConfig())
	err := svc.VerifyOTPAndResetPassword(context.Background(), "unknown@test.com", "123456", "newpassword2")

	if !errors.Is(err, ErrInvalidOTP) {
		t.Errorf("VerifyOTPAndResetPassword() error = %v, want ErrInvalidOTP (for security)", err)
	}
}

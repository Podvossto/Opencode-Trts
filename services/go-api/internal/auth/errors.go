// Purpose: Auth domain — error definitions
package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrUserNotFound       = errors.New("user not found")
	ErrSamePassword       = errors.New("new password must differ from current password")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrAdminNoOTP         = errors.New("admin accounts cannot use OTP password reset")
)

// Purpose: HTTP handlers for auth endpoints
// Ref: API Contract — POST /api/v1/auth/login, /logout, /password/change, /otp/request, /password/reset
// DB tables: users, otp_logs
package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tigersoft/ats-go-api/pkg/response"
)

// Handler holds auth HTTP handlers
type Handler struct {
	svc *Service
}

// NewHandler creates a new auth handler
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) || errors.Is(err, ErrAccountDisabled) {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
			return
		}
		log.Printf("login error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, resp)
}

// Logout handles POST /api/v1/auth/logout (requires JWT)
func (h *Handler) Logout(c *gin.Context) {
	jti, exists := c.Get("jti")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// Blacklist token for remaining TTL (30 minutes max)
	if err := h.svc.Logout(c.Request.Context(), jti.(string), 30*time.Minute); err != nil {
		log.Printf("logout error: %v", err)
		response.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// ChangePassword handles PUT /api/v1/auth/password/change (requires JWT)
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	err := h.svc.ChangePassword(c.Request.Context(), userID.(string), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "incorrect current password"})
			return
		}
		if errors.Is(err, ErrSamePassword) {
			response.BadRequest(c, []response.ValidationError{{Field: "new_password", Message: err.Error()}})
			return
		}
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c)
			return
		}
		log.Printf("change password error: %v", err)
		response.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// RequestOTP handles POST /api/v1/auth/otp/request
func (h *Handler) RequestOTP(c *gin.Context) {
	var req OTPRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	otp, err := h.svc.RequestOTP(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, ErrAdminNoOTP) {
			c.JSON(http.StatusForbidden, response.ErrorResponse{Error: err.Error()})
			return
		}
		log.Printf("request OTP error: %v", err)
		// Still return success to prevent email enumeration
	}

	// In development, log the OTP; in production, this would be sent via email
	if otp != "" {
		log.Printf("[DEV] OTP for %s: %s", req.Email, otp)
	}

	// Always return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{"message": "if the email exists, an OTP has been sent"})
}

// ResetPasswordOTP handles POST /api/v1/auth/password/reset
func (h *Handler) ResetPasswordOTP(c *gin.Context) {
	var req PasswordResetBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	err := h.svc.VerifyOTPAndResetPassword(c.Request.Context(), req.Email, req.OTP, req.NewPassword)
	if err != nil {
		if errors.Is(err, ErrInvalidOTP) {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
			return
		}
		log.Printf("reset password error: %v", err)
		response.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}

// Purpose: Users domain — shared types and request/response models
// DB table: users
package users

import "github.com/google/uuid"

// User is the domain model for a system user
type User struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	Role               string     `json:"role"`
	DepartmentID       *uuid.UUID `json:"department_id,omitempty"`
	IsActive           bool       `json:"is_active"`
	MustChangePassword bool       `json:"must_change_password"`
}

// CreateUserRequest is the body for POST /api/v1/admin/users
// New users get a system-generated temporary password; must_change_password defaults to true
type CreateUserRequest struct {
	Email        string     `json:"email" binding:"required,email"`
	FirstName    string     `json:"first_name" binding:"required"`
	LastName     string     `json:"last_name" binding:"required"`
	Role         string     `json:"role" binding:"required,oneof=admin hr supervisor"`
	DepartmentID *uuid.UUID `json:"department_id"`
}

// UpdateUserRequest is the body for PUT /api/v1/admin/users/:id
type UpdateUserRequest struct {
	FirstName    *string    `json:"first_name"`
	LastName     *string    `json:"last_name"`
	Role         *string    `json:"role" binding:"omitempty,oneof=admin hr supervisor"`
	DepartmentID *uuid.UUID `json:"department_id"`
	IsActive     *bool      `json:"is_active"`
}

// AdminResetPasswordRequest is the body for POST /api/v1/admin/users/:id/reset-password
type AdminResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

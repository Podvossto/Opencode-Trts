// Purpose: Users domain — service layer (business logic)
// Ref: API Contract — GET/POST /api/v1/admin/users, GET/PUT/DELETE /api/v1/admin/users/:id
// ATS-008: List users, ATS-009: Create user, ATS-010: Edit/Delete, ATS-011: Admin reset password
package users

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Domain errors
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrEmailExists    = errors.New("email already exists")
	ErrInvalidRole    = errors.New("invalid role")
	ErrSelfDeactivate = errors.New("cannot deactivate yourself")
	ErrDeptNotFound   = errors.New("department not found")
	ErrDeptExists     = errors.New("department name already exists")
)

// UpdateMeRequest is the body for PUT /api/v1/users/me (self-service)
type UpdateMeRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// CreateDepartmentRequest is the body for POST /api/v1/admin/departments
type CreateDepartmentRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

// Department is the response model
type Department struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// PaginatedResponse wraps a paginated list
type PaginatedResponse struct {
	Data       []User `json:"data"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}

// Service implements the user management business logic
type Service struct {
	repo RepositoryInterface
}

// NewService creates a new users service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// NewServiceWithDeps creates a new users service with explicit interface dep (for testing)
func NewServiceWithDeps(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

// List retrieves paginated users with filters (ATS-008)
func (s *Service) List(ctx context.Context, page, limit int, filter ListFilter) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rows, total, err := s.repo.List(ctx, page, limit, filter)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	users := make([]User, len(rows))
	for i, row := range rows {
		users[i] = rowToUser(row)
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetByID retrieves a single user
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if row == nil {
		return nil, ErrUserNotFound
	}
	u := rowToUser(*row)
	return &u, nil
}

// Create creates a new user with a temporary password (ATS-009)
// Admin only. New users get must_change_password=true.
func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, string, error) {
	// Validate role
	if req.Role != "admin" && req.Role != "hr" && req.Role != "supervisor" {
		return nil, "", ErrInvalidRole
	}

	// Check email uniqueness
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, "", ErrEmailExists
	}

	// Generate temporary password (12 chars alphanumeric)
	tempPwd, err := generateTempPassword(12)
	if err != nil {
		return nil, "", fmt.Errorf("generate temp password: %w", err)
	}

	row, err := s.repo.Create(ctx, req, tempPwd)
	if err != nil {
		return nil, "", fmt.Errorf("create user: %w", err)
	}

	u := rowToUser(*row)
	return &u, tempPwd, nil
}

// Update modifies a user's profile (ATS-010)
func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*User, error) {
	// Validate role if provided
	if req.Role != nil {
		if *req.Role != "admin" && *req.Role != "hr" && *req.Role != "supervisor" {
			return nil, ErrInvalidRole
		}
	}

	row, err := s.repo.Update(ctx, id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	if row == nil {
		return nil, ErrUserNotFound
	}
	u := rowToUser(*row)
	return &u, nil
}

// Delete soft-deletes (deactivates) a user (ATS-010)
func (s *Service) Delete(ctx context.Context, id uuid.UUID, callerID string) error {
	// Prevent self-deactivation
	callerUUID, err := uuid.Parse(callerID)
	if err == nil && callerUUID == id {
		return ErrSelfDeactivate
	}

	if err := s.repo.Deactivate(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// AdminResetPassword resets a user's password and sets must_change_password=true (ATS-011)
func (s *Service) AdminResetPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	if err := s.repo.ResetPassword(ctx, id, newPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("admin reset password: %w", err)
	}
	return nil
}

// GetMe retrieves the current user's profile
func (s *Service) GetMe(ctx context.Context, userID string) (*User, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return s.GetByID(ctx, uid)
}

// UpdateMe updates the current user's own profile (name only)
func (s *Service) UpdateMe(ctx context.Context, userID string, req UpdateMeRequest) (*User, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	row, err := s.repo.UpdateSelf(ctx, uid, req.FirstName, req.LastName)
	if err != nil {
		return nil, fmt.Errorf("update me: %w", err)
	}
	if row == nil {
		return nil, ErrUserNotFound
	}
	u := rowToUser(*row)
	return &u, nil
}

// ListDepartments returns all departments
func (s *Service) ListDepartments(ctx context.Context) ([]Department, error) {
	rows, err := s.repo.ListDepartments(ctx)
	if err != nil {
		return nil, fmt.Errorf("list departments: %w", err)
	}

	depts := make([]Department, len(rows))
	for i, row := range rows {
		depts[i] = Department{ID: row.ID, Name: row.Name}
	}
	return depts, nil
}

// CreateDepartment creates a new department
func (s *Service) CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*Department, error) {
	row, err := s.repo.CreateDepartment(ctx, req.Name)
	if err != nil {
		// Check for unique constraint violation (dept name already exists)
		if isDuplicateKeyError(err) {
			return nil, ErrDeptExists
		}
		return nil, fmt.Errorf("create department: %w", err)
	}
	return &Department{ID: row.ID, Name: row.Name}, nil
}

// -- helpers --

// rowToUser converts a UserRow to the domain User model
func rowToUser(row UserRow) User {
	return User{
		ID:                 row.ID,
		Email:              row.Email,
		FirstName:          row.FirstName,
		LastName:           row.LastName,
		Role:               row.Role,
		DepartmentID:       row.DepartmentID,
		IsActive:           row.IsActive,
		MustChangePassword: row.MustChangePassword,
	}
}

// generateTempPassword creates a random alphanumeric password of the given length
func generateTempPassword(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[n.Int64()]
	}
	return string(result), nil
}

// isDuplicateKeyError checks if a pgx error is a unique constraint violation
func isDuplicateKeyError(err error) bool {
	return err != nil && (fmt.Sprintf("%v", err) != "" &&
		(contains(err.Error(), "duplicate key") || contains(err.Error(), "23505")))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

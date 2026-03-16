// Purpose: Users domain — repository layer (database queries)
// DB tables: users, departments
package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// UserRow represents a full row from the users table
type UserRow struct {
	ID                 uuid.UUID
	Email              string
	PasswordHash       string
	FirstName          string
	LastName           string
	Role               string
	DepartmentID       *uuid.UUID
	DepartmentName     *string // joined from departments
	IsActive           bool
	MustChangePassword bool
}

// DepartmentRow represents a row from the departments table
type DepartmentRow struct {
	ID   uuid.UUID
	Name string
}

// ListFilter defines filters for listing users
type ListFilter struct {
	Search   string // search by name or email
	Role     string // filter by role
	DeptID   string // filter by department_id
	IsActive *bool  // filter by is_active
}

// RepositoryInterface defines the contract for users data access (enables mocking in tests)
type RepositoryInterface interface {
	List(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*UserRow, error)
	GetByEmail(ctx context.Context, email string) (*UserRow, error)
	Create(ctx context.Context, req CreateUserRequest, tempPassword string) (*UserRow, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserRow, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
	ResetPassword(ctx context.Context, id uuid.UUID, newPassword string) error
	UpdateSelf(ctx context.Context, id uuid.UUID, firstName, lastName string) (*UserRow, error)
	ListDepartments(ctx context.Context) ([]DepartmentRow, error)
	CreateDepartment(ctx context.Context, name string) (*DepartmentRow, error)
}

// Repository handles user-related database operations
type Repository struct {
	db *pgxpool.Pool
}

// compile-time check: Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

// NewRepository creates a new users repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// List retrieves paginated users with optional filters
// Returns (users, totalCount, error)
func (r *Repository) List(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error) {
	// Build WHERE clauses dynamically
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(u.first_name ILIKE $%d OR u.last_name ILIKE $%d OR u.email ILIKE $%d)",
			argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.Role != "" {
		conditions = append(conditions, fmt.Sprintf("u.role = $%d", argIdx))
		args = append(args, filter.Role)
		argIdx++
	}
	if filter.DeptID != "" {
		deptUUID, err := uuid.Parse(filter.DeptID)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("u.department_id = $%d", argIdx))
			args = append(args, deptUUID)
			argIdx++
		}
	}
	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("u.is_active = $%d", argIdx))
		args = append(args, *filter.IsActive)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users u %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Fetch page with LEFT JOIN to departments for department_name
	offset := (page - 1) * limit
	dataQuery := fmt.Sprintf(`
		SELECT u.id, u.email, u.first_name, u.last_name, u.role, u.department_id,
		       d.name AS department_name, u.is_active, u.must_change_password
		FROM users u
		LEFT JOIN departments d ON d.id = u.department_id
		%s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []UserRow
	for rows.Next() {
		var u UserRow
		if err := rows.Scan(
			&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role,
			&u.DepartmentID, &u.DepartmentName, &u.IsActive, &u.MustChangePassword,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user row: %w", err)
		}
		users = append(users, u)
	}

	return users, total, nil
}

// GetByID retrieves a single user by ID with department name
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*UserRow, error) {
	var u UserRow
	err := r.db.QueryRow(ctx, `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.role,
		       u.department_id, d.name AS department_name, u.is_active, u.must_change_password
		FROM users u
		LEFT JOIN departments d ON d.id = u.department_id
		WHERE u.id = $1`, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role,
			&u.DepartmentID, &u.DepartmentName, &u.IsActive, &u.MustChangePassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

// GetByEmail checks if a user with the given email exists (for uniqueness validation)
func (r *Repository) GetByEmail(ctx context.Context, email string) (*UserRow, error) {
	var u UserRow
	err := r.db.QueryRow(ctx,
		`SELECT id, email FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

// Create inserts a new user with a hashed temporary password
func (r *Repository) Create(ctx context.Context, req CreateUserRequest, tempPassword string) (*UserRow, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	var u UserRow
	err = r.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, first_name, last_name, role, department_id, must_change_password)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		RETURNING id, email, first_name, last_name, role, department_id, is_active, must_change_password`,
		req.Email, string(hash), req.FirstName, req.LastName, req.Role, req.DepartmentID).
		Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role,
			&u.DepartmentID, &u.IsActive, &u.MustChangePassword)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &u, nil
}

// Update modifies a user's profile fields (partial update)
func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserRow, error) {
	var sets []string
	var args []interface{}
	argIdx := 1

	if req.FirstName != nil {
		sets = append(sets, fmt.Sprintf("first_name = $%d", argIdx))
		args = append(args, *req.FirstName)
		argIdx++
	}
	if req.LastName != nil {
		sets = append(sets, fmt.Sprintf("last_name = $%d", argIdx))
		args = append(args, *req.LastName)
		argIdx++
	}
	if req.Role != nil {
		sets = append(sets, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, *req.Role)
		argIdx++
	}
	if req.DepartmentID != nil {
		sets = append(sets, fmt.Sprintf("department_id = $%d", argIdx))
		args = append(args, *req.DepartmentID)
		argIdx++
	}
	if req.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
		argIdx++
	}

	if len(sets) == 0 {
		return r.GetByID(ctx, id)
	}

	sets = append(sets, "updated_at = now()")
	query := fmt.Sprintf(`
		UPDATE users SET %s WHERE id = $%d
		RETURNING id, email, first_name, last_name, role, department_id, is_active, must_change_password`,
		strings.Join(sets, ", "), argIdx)
	args = append(args, id)

	var u UserRow
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role,
			&u.DepartmentID, &u.IsActive, &u.MustChangePassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &u, nil
}

// Deactivate sets is_active=false for a user (soft delete)
func (r *Repository) Deactivate(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `UPDATE users SET is_active = false, updated_at = now() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deactivate user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// ResetPassword sets a new password hash and must_change_password=true
func (r *Repository) ResetPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET password_hash = $1, must_change_password = true, updated_at = now() WHERE id = $2`,
		string(hash), id)
	if err != nil {
		return fmt.Errorf("reset password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// UpdateSelf updates only the fields a user can change about themselves (name only)
func (r *Repository) UpdateSelf(ctx context.Context, id uuid.UUID, firstName, lastName string) (*UserRow, error) {
	var u UserRow
	err := r.db.QueryRow(ctx, `
		UPDATE users SET first_name = $1, last_name = $2, updated_at = now()
		WHERE id = $3
		RETURNING id, email, first_name, last_name, role, department_id, is_active, must_change_password`,
		firstName, lastName, id).
		Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role,
			&u.DepartmentID, &u.IsActive, &u.MustChangePassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("update self: %w", err)
	}
	return &u, nil
}

// ListDepartments retrieves all departments ordered by name
func (r *Repository) ListDepartments(ctx context.Context) ([]DepartmentRow, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM departments ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list departments: %w", err)
	}
	defer rows.Close()

	var depts []DepartmentRow
	for rows.Next() {
		var d DepartmentRow
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, fmt.Errorf("scan department: %w", err)
		}
		depts = append(depts, d)
	}
	return depts, nil
}

// CreateDepartment inserts a new department
func (r *Repository) CreateDepartment(ctx context.Context, name string) (*DepartmentRow, error) {
	var d DepartmentRow
	err := r.db.QueryRow(ctx,
		`INSERT INTO departments (name) VALUES ($1) RETURNING id, name`, name).
		Scan(&d.ID, &d.Name)
	if err != nil {
		return nil, fmt.Errorf("create department: %w", err)
	}
	return &d, nil
}

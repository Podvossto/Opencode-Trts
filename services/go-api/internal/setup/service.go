// Purpose: Business logic for Setup — reference check before delete (blocks if in-use),
//
//	SYS template read-only enforcement, async JD embedding dispatch to ai-service.
//
// Ref: US-M14-01 through US-M14-04, API Contract — /api/v1/setup/*
// Dev must implement: SetupService interface — CRUD per entity, DispatchEmbedding
// DB tables: departments, email_templates, system_settings
package setup

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrDepartmentNotFound is returned when a department does not exist
var ErrDepartmentNotFound = errors.New("department not found")

// ErrDepartmentInUse is returned when deleting a department with active jobs
var ErrDepartmentInUse = errors.New("department has active jobs and cannot be deleted")

// ErrTemplateNotFound is returned when an email template slug does not match any row
var ErrTemplateNotFound = errors.New("email template not found")

// ErrSystemTemplateReadOnly is returned when trying to delete or hard-modify a SYS_ template
var ErrSystemTemplateReadOnly = errors.New("system email templates cannot be deleted")

// setupService is the concrete implementation of SetupService
// Dev must implement: inject pgxpool.Pool
type setupService struct {
	// db *pgxpool.Pool
}

// NewService constructs a production SetupService
// Dev must implement: accept *pgxpool.Pool as parameter
func NewService() SetupService {
	return &setupService{}
}

// ListDepartments retrieves all departments ordered by name
// Dev must implement: SELECT * FROM departments ORDER BY name ASC
func (s *setupService) ListDepartments(ctx context.Context) ([]Department, error) {
	return nil, errors.New("not implemented")
}

// CreateDepartment inserts a new department
// Guard: code must be unique across departments
// Dev must implement:
//
//	INSERT INTO departments (id, name, code, manager_id) VALUES (gen_random_uuid(),$1,$2,$3) RETURNING *
func (s *setupService) CreateDepartment(ctx context.Context, req UpsertDepartmentRequest) (*Department, error) {
	return nil, errors.New("not implemented")
}

// UpdateDepartment modifies an existing department
// Dev must implement:
//
//	UPDATE departments SET name=$1, code=$2, manager_id=$3 WHERE id=$4 RETURNING *
func (s *setupService) UpdateDepartment(ctx context.Context, id uuid.UUID, req UpsertDepartmentRequest) (*Department, error) {
	return nil, errors.New("not implemented")
}

// DeleteDepartment removes a department if it has no active job postings
// Dev must implement:
//
//	SELECT COUNT(*) FROM jobs WHERE department_id=$1 AND status IN ('open','in_progress')
//	If count > 0: return ErrDepartmentInUse
//	DELETE FROM departments WHERE id=$1
func (s *setupService) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}

// ListEmailTemplates retrieves all email templates ordered by slug
// Dev must implement: SELECT * FROM email_templates ORDER BY slug ASC
func (s *setupService) ListEmailTemplates(ctx context.Context) ([]EmailTemplate, error) {
	return nil, errors.New("not implemented")
}

// GetEmailTemplate retrieves a single template by slug
// Dev must implement: SELECT * FROM email_templates WHERE slug=$1
func (s *setupService) GetEmailTemplate(ctx context.Context, slug string) (*EmailTemplate, error) {
	// TODO: return ErrTemplateNotFound if no rows
	return nil, errors.New("not implemented")
}

// UpdateEmailTemplate modifies subject and body of a template
// Guard: SYS_ slugs (system templates) cannot have their variables array changed
// Dev must implement:
//
//	UPDATE email_templates SET subject=$1, body=$2, updated_at=now() WHERE slug=$3 RETURNING *
func (s *setupService) UpdateEmailTemplate(ctx context.Context, slug string, req UpdateEmailTemplateRequest) (*EmailTemplate, error) {
	return nil, errors.New("not implemented")
}

// GetSettings retrieves all system settings as key-value pairs
// Dev must implement: SELECT key, value, updated_at FROM system_settings ORDER BY key
func (s *setupService) GetSettings(ctx context.Context) ([]SystemSetting, error) {
	return nil, errors.New("not implemented")
}

// UpdateSetting upserts a system setting by key
// Dev must implement:
//
//	INSERT INTO system_settings (key, value) VALUES ($1,$2)
//	ON CONFLICT (key) DO UPDATE SET value=$2, updated_at=now() RETURNING *
func (s *setupService) UpdateSetting(ctx context.Context, key string, value string) (*SystemSetting, error) {
	return nil, errors.New("not implemented")
}

// Purpose: Jobs domain — repository layer (database queries)
// DB tables: jobs, departments, employment_types, experience_levels, application_forms
// ATS-012: Portal public job listing + detail
package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PublicJobRow represents a job row enriched with joined names for public display
type PublicJobRow struct {
	ID                  uuid.UUID
	Title               string
	DepartmentName      string
	EmploymentTypeName  string
	ExperienceLevelName string
	Description         string
	Requirements        *string
	Slots               int
	PublishAt           *time.Time
	CloseAt             *time.Time
	CreatedAt           time.Time
}

// PublicJobDetailRow extends PublicJobRow with form_schema for the detail endpoint
type PublicJobDetailRow struct {
	PublicJobRow
	FormID     uuid.UUID
	FormSchema string // JSONB as string
}

// Repository handles job-related database operations
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new jobs repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ListPublicJobs retrieves all open jobs with pagination for the portal
// Only shows jobs with status='open'
func (r *Repository) ListPublicJobs(ctx context.Context, page, limit int, search string) ([]PublicJobRow, int, error) {
	baseWhere := "WHERE j.status = 'open'"
	var args []interface{}
	argIdx := 1

	if search != "" {
		baseWhere += fmt.Sprintf(" AND (j.title ILIKE $%d OR j.description ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs j %s", baseWhere)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count public jobs: %w", err)
	}

	// Fetch page
	offset := (page - 1) * limit
	dataQuery := fmt.Sprintf(`
		SELECT j.id, j.title, d.name AS department_name,
		       et.name AS employment_type_name, el.name AS experience_level_name,
		       j.description, j.requirements, j.slots, j.publish_at, j.close_at, j.created_at
		FROM jobs j
		JOIN departments d ON d.id = j.department_id
		JOIN employment_types et ON et.id = j.employment_type_id
		JOIN experience_levels el ON el.id = j.experience_level_id
		%s
		ORDER BY j.publish_at DESC NULLS LAST, j.created_at DESC
		LIMIT $%d OFFSET $%d`,
		baseWhere, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list public jobs: %w", err)
	}
	defer rows.Close()

	var jobs []PublicJobRow
	for rows.Next() {
		var j PublicJobRow
		if err := rows.Scan(
			&j.ID, &j.Title, &j.DepartmentName, &j.EmploymentTypeName,
			&j.ExperienceLevelName, &j.Description, &j.Requirements,
			&j.Slots, &j.PublishAt, &j.CloseAt, &j.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan public job row: %w", err)
		}
		jobs = append(jobs, j)
	}

	return jobs, total, nil
}

// GetPublicJobByID retrieves a single open job with its application form schema
func (r *Repository) GetPublicJobByID(ctx context.Context, id uuid.UUID) (*PublicJobDetailRow, error) {
	var j PublicJobDetailRow
	err := r.db.QueryRow(ctx, `
		SELECT j.id, j.title, d.name AS department_name,
		       et.name AS employment_type_name, el.name AS experience_level_name,
		       j.description, j.requirements, j.slots, j.publish_at, j.close_at, j.created_at,
		       j.form_id, COALESCE(af.form_schema::TEXT, '[]') AS form_schema
		FROM jobs j
		JOIN departments d ON d.id = j.department_id
		JOIN employment_types et ON et.id = j.employment_type_id
		JOIN experience_levels el ON el.id = j.experience_level_id
		LEFT JOIN application_forms af ON af.id = j.form_id
		WHERE j.id = $1 AND j.status = 'open'`, id).
		Scan(&j.ID, &j.Title, &j.DepartmentName, &j.EmploymentTypeName,
			&j.ExperienceLevelName, &j.Description, &j.Requirements,
			&j.Slots, &j.PublishAt, &j.CloseAt, &j.CreatedAt,
			&j.FormID, &j.FormSchema)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get public job by id: %w", err)
	}
	return &j, nil
}

// DashboardStats holds the pipeline dashboard statistics for a job
type DashboardStats struct {
	Total     int `json:"total"`
	Remaining int `json:"remaining"`
	Dropped   int `json:"dropped"`
	Hired     int `json:"hired"`
}

// GetDashboardStats retrieves pipeline stats for a job
func (r *Repository) GetDashboardStats(ctx context.Context, jobID uuid.UUID) (*DashboardStats, error) {
	var stats DashboardStats

	err := r.db.QueryRow(ctx, `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status NOT IN ('dropped', 'hired', 'transferred')) as remaining,
			COUNT(*) FILTER (WHERE status = 'dropped') as dropped,
			COUNT(*) FILTER (WHERE status = 'hired') as hired
		FROM applications WHERE job_id = $1`, jobID).
		Scan(&stats.Total, &stats.Remaining, &stats.Dropped, &stats.Hired)
	if err != nil {
		return nil, fmt.Errorf("get dashboard stats: %w", err)
	}

	return &stats, nil
}

// JobExists checks if a job exists by ID
func (r *Repository) JobExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM jobs WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check job exists: %w", err)
	}
	return exists, nil
}

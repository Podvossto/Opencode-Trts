// Purpose: Jobs domain — repository layer (database queries)
// DB tables: jobs, departments, employment_types, experience_levels, application_forms, job_hard_skills
// ATS-012: Portal public job listing + detail
// ATS-026: Job CRUD with hard skill weights
package jobs

import (
	"context"
	"fmt"
	"strings"
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

// CreateJob inserts a new job with its hard skills in a transaction
func (r *Repository) CreateJob(ctx context.Context, job Job, skills []JobHardSkill) (Job, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Job{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var id uuid.UUID
	var createdAt time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO jobs (title, department_id, employment_type_id, experience_level_id, 
		                  form_id, description, requirements, slots, status, publish_at, close_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`,
		job.Title, job.DepartmentID, job.EmploymentTypeID, job.ExperienceLevelID,
		job.FormID, job.Description, job.Requirements, job.Slots, job.Status,
		job.PublishAt, job.CloseAt, job.CreatedBy).
		Scan(&id, &createdAt)
	if err != nil {
		return Job{}, fmt.Errorf("insert job: %w", err)
	}

	for _, skill := range skills {
		_, err = tx.Exec(ctx, `
			INSERT INTO job_hard_skills (job_id, skill_name, weight)
			VALUES ($1, $2, $3)`,
			id, skill.SkillName, skill.Weight)
		if err != nil {
			return Job{}, fmt.Errorf("insert job hard skill: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return Job{}, fmt.Errorf("commit transaction: %w", err)
	}

	job.ID = id
	job.CreatedAt = createdAt
	return job, nil
}

// ListJobs retrieves jobs with filtering and pagination
func (r *Repository) ListJobs(ctx context.Context, filter JobFilter) ([]Job, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("j.status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Department != "" {
		conditions = append(conditions, fmt.Sprintf("d.name ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Department+"%")
		argIdx++
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(j.title ILIKE $%d OR j.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs j JOIN departments d ON d.id = j.department_id %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count jobs: %w", err)
	}

	offset := (filter.Page - 1) * filter.Limit
	dataQuery := fmt.Sprintf(`
		SELECT j.id, j.title, j.department_id, j.employment_type_id, j.experience_level_id,
		       j.form_id, j.description, j.requirements, j.slots, j.status, 
		       j.publish_at, j.close_at, j.created_at, j.created_by
		FROM jobs j
		JOIN departments d ON d.id = j.department_id
		%s
		ORDER BY j.created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIdx, argIdx+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(
			&j.ID, &j.Title, &j.DepartmentID, &j.EmploymentTypeID, &j.ExperienceLevelID,
			&j.FormID, &j.Description, &j.Requirements, &j.Slots, &j.Status,
			&j.PublishAt, &j.CloseAt, &j.CreatedAt, &j.CreatedBy,
		); err != nil {
			return nil, 0, fmt.Errorf("scan job row: %w", err)
		}
		jobs = append(jobs, j)
	}

	return jobs, total, nil
}

// GetJobByID retrieves a job by ID with its hard skills
func (r *Repository) GetJobByID(ctx context.Context, id string) (Job, []JobHardSkill, error) {
	var job Job
	var parsedID uuid.UUID
	var err error

	parsedID, err = uuid.Parse(id)
	if err != nil {
		return job, nil, fmt.Errorf("invalid UUID: %w", err)
	}

	err = r.db.QueryRow(ctx, `
		SELECT id, title, department_id, employment_type_id, experience_level_id,
		       form_id, description, requirements, slots, status, 
		       publish_at, close_at, created_at, created_by
		FROM jobs WHERE id = $1`, parsedID).
		Scan(&job.ID, &job.Title, &job.DepartmentID, &job.EmploymentTypeID, &job.ExperienceLevelID,
			&job.FormID, &job.Description, &job.Requirements, &job.Slots, &job.Status,
			&job.PublishAt, &job.CloseAt, &job.CreatedAt, &job.CreatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return job, nil, nil
		}
		return job, nil, fmt.Errorf("get job by id: %w", err)
	}

	skillRows, err := r.db.Query(ctx, `
		SELECT id, job_id, skill_name, weight
		FROM job_hard_skills WHERE job_id = $1`, parsedID)
	if err != nil {
		return job, nil, fmt.Errorf("get job hard skills: %w", err)
	}
	defer skillRows.Close()

	var skills []JobHardSkill
	for skillRows.Next() {
		var skill JobHardSkill
		if err := skillRows.Scan(&skill.ID, &skill.JobID, &skill.SkillName, &skill.Weight); err != nil {
			return job, nil, fmt.Errorf("scan hard skill: %w", err)
		}
		skills = append(skills, skill)
	}

	return job, skills, nil
}

// UpdateJob updates a job and its hard skills in a transaction
func (r *Repository) UpdateJob(ctx context.Context, id string, job Job, skills []JobHardSkill) (Job, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Job{}, fmt.Errorf("invalid UUID: %w", err)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Job{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE jobs SET title = $1, department_id = $2, employment_type_id = $3, 
		       experience_level_id = $4, form_id = $5, description = $6, 
		       requirements = $7, slots = $8, status = $9, publish_at = $10, 
		       close_at = $11, updated_at = NOW()
		WHERE id = $12`,
		job.Title, job.DepartmentID, job.EmploymentTypeID, job.ExperienceLevelID,
		job.FormID, job.Description, job.Requirements, job.Slots, job.Status,
		job.PublishAt, job.CloseAt, parsedID)
	if err != nil {
		return Job{}, fmt.Errorf("update job: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM job_hard_skills WHERE job_id = $1`, parsedID)
	if err != nil {
		return Job{}, fmt.Errorf("delete old skills: %w", err)
	}

	for _, skill := range skills {
		_, err = tx.Exec(ctx, `
			INSERT INTO job_hard_skills (job_id, skill_name, weight)
			VALUES ($1, $2, $3)`,
			parsedID, skill.SkillName, skill.Weight)
		if err != nil {
			return Job{}, fmt.Errorf("insert job hard skill: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return Job{}, fmt.Errorf("commit transaction: %w", err)
	}

	job.ID = parsedID
	return job, nil
}

// DeleteJob deletes a job by ID (only if not open)
func (r *Repository) DeleteJob(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	var status string
	err = r.db.QueryRow(ctx, `SELECT status FROM jobs WHERE id = $1`, parsedID).Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("job not found")
		}
		return fmt.Errorf("get job status: %w", err)
	}

	if status == "open" {
		return fmt.Errorf("cannot delete open job")
	}

	_, err = r.db.Exec(ctx, `DELETE FROM jobs WHERE id = $1`, parsedID)
	if err != nil {
		return fmt.Errorf("delete job: %w", err)
	}

	return nil
}

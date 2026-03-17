// Purpose: Jobs domain — service layer (business logic)
// ATS-012: Portal public job listing + detail
// ATS-026: Job CRUD with hard skill weights
package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PublicJob is the response model for public job listings
type PublicJob struct {
	ID                  uuid.UUID  `json:"id"`
	Title               string     `json:"title"`
	DepartmentName      string     `json:"department_name"`
	EmploymentTypeName  string     `json:"employment_type_name"`
	ExperienceLevelName string     `json:"experience_level_name"`
	Description         string     `json:"description"`
	Requirements        *string    `json:"requirements,omitempty"`
	Slots               int        `json:"slots"`
	PublishAt           *time.Time `json:"publish_at,omitempty"`
	CloseAt             *time.Time `json:"close_at,omitempty"`
}

// PublicJobDetail extends PublicJob with application form schema
type PublicJobDetail struct {
	PublicJob
	FormSchema json.RawMessage `json:"form_schema"`
}

// PublicJobListResponse wraps paginated public jobs
type PublicJobListResponse struct {
	Data       []PublicJob `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// PortalService implements portal-specific business logic
type PortalService struct {
	repo *Repository
}

// NewPortalService creates a new portal service
func NewPortalService(repo *Repository) *PortalService {
	return &PortalService{repo: repo}
}

// ListPublicJobs retrieves paginated open jobs for the portal (ATS-012)
func (s *PortalService) ListPublicJobs(ctx context.Context, page, limit int, search string) (*PublicJobListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	rows, total, err := s.repo.ListPublicJobs(ctx, page, limit, search)
	if err != nil {
		return nil, fmt.Errorf("list public jobs: %w", err)
	}

	jobs := make([]PublicJob, len(rows))
	for i, row := range rows {
		jobs[i] = PublicJob{
			ID:                  row.ID,
			Title:               row.Title,
			DepartmentName:      row.DepartmentName,
			EmploymentTypeName:  row.EmploymentTypeName,
			ExperienceLevelName: row.ExperienceLevelName,
			Description:         row.Description,
			Requirements:        row.Requirements,
			Slots:               row.Slots,
			PublishAt:           row.PublishAt,
			CloseAt:             row.CloseAt,
		}
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return &PublicJobListResponse{
		Data:       jobs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetPublicJob retrieves a single open job with its application form (ATS-012)
func (s *PortalService) GetPublicJob(ctx context.Context, id uuid.UUID) (*PublicJobDetail, error) {
	row, err := s.repo.GetPublicJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get public job: %w", err)
	}
	if row == nil {
		return nil, nil // not found or not open
	}

	return &PublicJobDetail{
		PublicJob: PublicJob{
			ID:                  row.ID,
			Title:               row.Title,
			DepartmentName:      row.DepartmentName,
			EmploymentTypeName:  row.EmploymentTypeName,
			ExperienceLevelName: row.ExperienceLevelName,
			Description:         row.Description,
			Requirements:        row.Requirements,
			Slots:               row.Slots,
			PublishAt:           row.PublishAt,
			CloseAt:             row.CloseAt,
		},
		FormSchema: json.RawMessage(row.FormSchema),
	}, nil
}

// GetDashboardStats retrieves pipeline statistics for a job (ATS-028)
func (s *PortalService) GetDashboardStats(ctx context.Context, jobID uuid.UUID) (*DashboardStats, error) {
	exists, err := s.repo.JobExists(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("check job exists: %w", err)
	}
	if !exists {
		return nil, nil
	}

	stats, err := s.repo.GetDashboardStats(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get dashboard stats: %w", err)
	}

	return stats, nil
}

// JobService implements job CRUD business logic
type JobService struct {
	repo *Repository
}

// NewJobService creates a new job service
func NewJobService(repo *Repository) *JobService {
	return &JobService{repo: repo}
}

// JobWithSkills represents a job with its hard skills
type JobWithSkills struct {
	Job        Job            `json:"job"`
	HardSkills []JobHardSkill `json:"hard_skills"`
}

// JobListResponse wraps paginated jobs
type JobListResponse struct {
	Data       []JobWithSkills `json:"data"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// CreateJob creates a new job with hard skills
func (s *JobService) CreateJob(ctx context.Context, job Job, skills []JobHardSkill) (Job, error) {
	if job.Status == "" {
		job.Status = "draft"
	}
	createdJob, err := s.repo.CreateJob(ctx, job, skills)
	if err != nil {
		return Job{}, fmt.Errorf("create job: %w", err)
	}
	return createdJob, nil
}

// ListJobs retrieves paginated jobs with filters
func (s *JobService) ListJobs(ctx context.Context, filter JobFilter) (*JobListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 50 {
		filter.Limit = 10
	}

	jobs, total, err := s.repo.ListJobs(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}

	data := make([]JobWithSkills, len(jobs))
	for i, job := range jobs {
		_, skills, err := s.repo.GetJobByID(ctx, job.ID.String())
		if err != nil {
			return nil, fmt.Errorf("get job skills: %w", err)
		}
		data[i] = JobWithSkills{
			Job:        job,
			HardSkills: skills,
		}
	}

	totalPages := total / filter.Limit
	if total%filter.Limit != 0 {
		totalPages++
	}

	return &JobListResponse{
		Data:       data,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetJobByID retrieves a job by ID with its hard skills
func (s *JobService) GetJobByID(ctx context.Context, id string) (*JobWithSkills, error) {
	job, skills, err := s.repo.GetJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job: %w", err)
	}
	if job.ID == uuid.Nil {
		return nil, nil
	}

	return &JobWithSkills{
		Job:        job,
		HardSkills: skills,
	}, nil
}

// UpdateJob updates a job and its hard skills
func (s *JobService) UpdateJob(ctx context.Context, id string, job Job, skills []JobHardSkill) (Job, error) {
	updatedJob, err := s.repo.UpdateJob(ctx, id, job, skills)
	if err != nil {
		return Job{}, fmt.Errorf("update job: %w", err)
	}
	return updatedJob, nil
}

// DeleteJob deletes a job by ID
func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	err := s.repo.DeleteJob(ctx, id)
	if err != nil {
		return fmt.Errorf("delete job: %w", err)
	}
	return nil
}

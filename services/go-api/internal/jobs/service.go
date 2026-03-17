// Purpose: Jobs domain — service layer (business logic)
// ATS-012: Portal public job listing + detail
// Full job CRUD (Sprint 2) will be added later
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

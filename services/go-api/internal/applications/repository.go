// Purpose: Applications domain — repository layer (database queries)
// DB tables: candidates, jobs, applications, application_documents
package applications

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CandidateRow represents a row from the candidates table
type CandidateRow struct {
	ID              uuid.UUID
	Email           string
	FirstName       string
	LastName        string
	Phone           *string
	IsBlacklisted   bool
	BlacklistReason *string
}

// ApplicationRow represents a row from the applications table
type ApplicationRow struct {
	ID              uuid.UUID
	CandidateID     uuid.UUID
	JobID           uuid.UUID
	Status          string
	FormData        map[string]interface{}
	OCRStatus       string
	SimilarityScore *float64
	SubmittedAt     pgtype.Timestamptz
	CreatedAt       pgtype.Timestamptz
}

// ApplicationDocumentRow represents a row from the application_documents table
type ApplicationDocumentRow struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	FileType      string
	FileURL       string
	FileName      string
	FileSizeKB    *int32
	OCRStatus     *string
	OCRText       *string
}

// JobRow represents a row from the jobs table (for validation)
type JobRow struct {
	ID           uuid.UUID
	Title        string
	Status       string
	DepartmentID uuid.UUID
}

// RepositoryInterface defines the contract for applications data access
type RepositoryInterface interface {
	// Candidate operations
	FindCandidateByEmail(ctx context.Context, email string) (*CandidateRow, error)
	CreateCandidate(ctx context.Context, email, firstName, lastName, phone string) (*CandidateRow, error)

	// Job operations
	GetJobByID(ctx context.Context, id uuid.UUID) (*JobRow, error)

	// Application operations
	CreateApplication(ctx context.Context, candidateID, jobID uuid.UUID, formData map[string]interface{}) (*ApplicationRow, error)
	GetApplicationByID(ctx context.Context, id uuid.UUID) (*ApplicationRow, error)
	GetApplicationByCandidateAndJob(ctx context.Context, candidateID, jobID uuid.UUID) (*ApplicationRow, error)

	// Document operations
	CreateApplicationDocument(ctx context.Context, appID uuid.UUID, fileType, fileURL, fileName string, fileSizeKB *int32) (*ApplicationDocumentRow, error)
	UpdateApplicationDocumentOCRStatus(ctx context.Context, docID uuid.UUID, status string, ocrText *string) error
}

// Repository handles application-related database operations
type Repository struct {
	db *pgxpool.Pool
}

// compile-time check: Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

// NewRepository creates a new applications repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// FindCandidateByEmail retrieves a candidate by email address
func (r *Repository) FindCandidateByEmail(ctx context.Context, email string) (*CandidateRow, error) {
	var c CandidateRow
	err := r.db.QueryRow(ctx, `
		SELECT id, email, first_name, last_name, phone, is_blacklisted, blacklist_reason
		FROM candidates WHERE email = $1`, email).
		Scan(&c.ID, &c.Email, &c.FirstName, &c.LastName, &c.Phone, &c.IsBlacklisted, &c.BlacklistReason)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find candidate by email: %w", err)
	}
	return &c, nil
}

// CreateCandidate inserts a new candidate
func (r *Repository) CreateCandidate(ctx context.Context, email, firstName, lastName, phone string) (*CandidateRow, error) {
	var c CandidateRow
	err := r.db.QueryRow(ctx, `
		INSERT INTO candidates (email, first_name, last_name, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, first_name, last_name, phone, is_blacklisted, blacklist_reason`,
		email, firstName, lastName, phone).
		Scan(&c.ID, &c.Email, &c.FirstName, &c.LastName, &c.Phone, &c.IsBlacklisted, &c.BlacklistReason)
	if err != nil {
		return nil, fmt.Errorf("create candidate: %w", err)
	}
	return &c, nil
}

// GetJobByID retrieves a job by ID
func (r *Repository) GetJobByID(ctx context.Context, id uuid.UUID) (*JobRow, error) {
	var j JobRow
	err := r.db.QueryRow(ctx, `
		SELECT id, title, status, department_id
		FROM jobs WHERE id = $1`, id).
		Scan(&j.ID, &j.Title, &j.Status, &j.DepartmentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get job by id: %w", err)
	}
	return &j, nil
}

// CreateApplication inserts a new application
func (r *Repository) CreateApplication(ctx context.Context, candidateID, jobID uuid.UUID, formData map[string]interface{}) (*ApplicationRow, error) {
	var a ApplicationRow
	err := r.db.QueryRow(ctx, `
		INSERT INTO applications (candidate_id, job_id, form_data, status, ocr_status)
		VALUES ($1, $2, $3, 'pending', 'pending')
		RETURNING id, candidate_id, job_id, status, form_data, ocr_status, similarity_score, submitted_at, created_at`,
		candidateID, jobID, formData).
		Scan(&a.ID, &a.CandidateID, &a.JobID, &a.Status, &a.FormData, &a.OCRStatus, &a.SimilarityScore, &a.SubmittedAt, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create application: %w", err)
	}
	return &a, nil
}

// GetApplicationByID retrieves an application by ID
func (r *Repository) GetApplicationByID(ctx context.Context, id uuid.UUID) (*ApplicationRow, error) {
	var a ApplicationRow
	err := r.db.QueryRow(ctx, `
		SELECT id, candidate_id, job_id, status, form_data, ocr_status, similarity_score, submitted_at, created_at
		FROM applications WHERE id = $1`, id).
		Scan(&a.ID, &a.CandidateID, &a.JobID, &a.Status, &a.FormData, &a.OCRStatus, &a.SimilarityScore, &a.SubmittedAt, &a.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get application by id: %w", err)
	}
	return &a, nil
}

// GetApplicationByCandidateAndJob checks if candidate already applied to this job
func (r *Repository) GetApplicationByCandidateAndJob(ctx context.Context, candidateID, jobID uuid.UUID) (*ApplicationRow, error) {
	var a ApplicationRow
	err := r.db.QueryRow(ctx, `
		SELECT id, candidate_id, job_id, status, form_data, ocr_status, similarity_score, submitted_at, created_at
		FROM applications WHERE candidate_id = $1 AND job_id = $2`, candidateID, jobID).
		Scan(&a.ID, &a.CandidateID, &a.JobID, &a.Status, &a.FormData, &a.OCRStatus, &a.SimilarityScore, &a.SubmittedAt, &a.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get application by candidate and job: %w", err)
	}
	return &a, nil
}

// CreateApplicationDocument inserts a new application document
func (r *Repository) CreateApplicationDocument(ctx context.Context, appID uuid.UUID, fileType, fileURL, fileName string, fileSizeKB *int32) (*ApplicationDocumentRow, error) {
	var d ApplicationDocumentRow
	err := r.db.QueryRow(ctx, `
		INSERT INTO application_documents (application_id, file_type, file_url, file_name, file_size_kb)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, application_id, file_type, file_url, file_name, file_size_kb`,
		appID, fileType, fileURL, fileName, fileSizeKB).
		Scan(&d.ID, &d.ApplicationID, &d.FileType, &d.FileURL, &d.FileName, &d.FileSizeKB)
	if err != nil {
		return nil, fmt.Errorf("create application document: %w", err)
	}
	return &d, nil
}

// UpdateApplicationDocumentOCRStatus updates the OCR status and text for a document
func (r *Repository) UpdateApplicationDocumentOCRStatus(ctx context.Context, docID uuid.UUID, status string, ocrText *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE application_documents
		SET ocr_status = $1, ocr_text = $2
		WHERE id = $3`, status, ocrText, docID)
	if err != nil {
		return fmt.Errorf("update document ocr status: %w", err)
	}
	return nil
}

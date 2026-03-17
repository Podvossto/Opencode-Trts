// Purpose: Applications domain — portal submission, autosave, status tracking
// Ref: API Contract — POST /api/v1/portal/applications, GET /api/v1/portal/applications/{id}
//
//	POST /api/v1/portal/applications/autosave
//	GET /api/v1/applications (HR pipeline view)
//
// DB tables: candidates, applications, application_documents, test_autosave
// Dev must implement: ApplicationHandler, ApplicationService, ApplicationRepository
// IMPORTANT: submission returns 201 immediately; OCR is async via goroutine worker
package applications

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
	"github.com/tigersoft/ats-go-api/internal/middleware"
	"github.com/tigersoft/ats-go-api/internal/ocr_worker"
)

// RegisterRoutes mounts all application routes.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	repo := NewRepository(db)
	ocrWorker := ocr_worker.NewWorker(db, cfg.AIServiceURL, 100)
	svc := NewApplicationService(repo, ocrWorker)
	h := NewApplicationHandler(svc)

	portal := rg.Group("/portal")
	{
		portal.POST("/apply", h.SubmitApplication)
		portal.GET("/applications/:id/status", h.GetApplicationStatus)
		portal.POST("/applications/:id/autosave", h.AutosaveApplication)
	}

	apps := rg.Group("/applications")
	apps.Use(middleware.AuthMiddleware(cfg, rdb))
	apps.Use(middleware.RequireRole("admin", "hr"))
	{
		apps.GET("", h.ListApplications)
		apps.GET("/:id", h.GetApplication)
	}
}

// ApplicationStatus is the pipeline stage enum
type ApplicationStatus string

const (
	StatusPending      ApplicationStatus = "pending"
	StatusInReview     ApplicationStatus = "in_review"
	StatusTesting      ApplicationStatus = "testing"
	StatusInterviewing ApplicationStatus = "interviewing"
	StatusHired        ApplicationStatus = "hired"
	StatusDropped      ApplicationStatus = "dropped"
	StatusTransferred  ApplicationStatus = "transferred"
)

// Application is the domain model
type Application struct {
	ID          uuid.UUID         `json:"id"`
	CandidateID uuid.UUID         `json:"candidate_id"`
	JobID       uuid.UUID         `json:"job_id"`
	Status      ApplicationStatus `json:"status"`
	FormData    interface{}       `json:"form_data"` // JSONB answers
	SubmittedAt time.Time         `json:"submitted_at"`
	CreatedAt   time.Time         `json:"created_at"`
}

// SubmitApplicationRequest is the portal submission body
type SubmitApplicationRequest struct {
	JobID     uuid.UUID              `json:"job_id" binding:"required"`
	FirstName string                 `json:"first_name" binding:"required"`
	LastName  string                 `json:"last_name" binding:"required"`
	Email     string                 `json:"email" binding:"required,email"`
	Phone     string                 `json:"phone" binding:"required"`
	FormData  map[string]interface{} `json:"form_data" binding:"required"`
}

// ApplicationService defines the application business logic contract.
type ApplicationService interface {
	Submit(ctx context.Context, req SubmitApplicationRequest) (*Application, error)
	SubmitWithDocument(ctx context.Context, req SubmitApplicationRequest, fileType, fileURL, fileName string, fileSizeKB *int32) (*Application, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Application, error)
	GetStatus(ctx context.Context, id uuid.UUID) (*ApplicationStatusResponse, error)
	Autosave(ctx context.Context, jobID uuid.UUID, email string, data map[string]interface{}) error
	ListForPipeline(ctx context.Context, jobID *uuid.UUID, statusFilter *string, page, limit int) ([]Application, int, error)
}

// ApplicationStatusResponse is the status view response
type ApplicationStatusResponse struct {
	ID     uuid.UUID         `json:"application_id"`
	Status ApplicationStatus `json:"status"`
	Stage  string            `json:"stage"`
}

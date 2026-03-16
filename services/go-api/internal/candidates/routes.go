// Purpose: Routes, models, and service interface for Candidates domain — talent pool management
// Ref: API Contract — GET /api/v1/candidates, GET /api/v1/candidates/{id}
//
//	PUT /api/v1/candidates/{id}/star, GET /api/v1/candidates/{id}/applications
//
// DB tables: candidates, applications
// Dev must implement: CandidateHandler, CandidateService, CandidateRepository
package candidates

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
)

// RegisterRoutes mounts all candidate management routes.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	// Dev must implement:
	// GET  /candidates             (HR/Supervisor)
	// GET  /candidates/:id         (HR/Supervisor)
	// PUT  /candidates/:id/star    (HR)
	// GET  /candidates/:id/history (HR/Supervisor) — all applications by this candidate
}

// Candidate is the domain model for a job applicant profile.
type Candidate struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	IsStarred bool      `json:"is_starred"`
	CreatedAt time.Time `json:"created_at"`
}

// CandidateService defines the candidate business logic contract.
type CandidateService interface {
	List(ctx context.Context, filters map[string]string) ([]Candidate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Candidate, error)
	SetStar(ctx context.Context, id uuid.UUID, starred bool) error
	GetApplicationHistory(ctx context.Context, id uuid.UUID) ([]interface{}, error)
}

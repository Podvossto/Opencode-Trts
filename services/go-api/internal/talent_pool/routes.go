// Purpose: Register HTTP routes for the talent pool domain (transferred/saved candidates)
// Ref: API Contract — GET/POST/DELETE /api/v1/talent-pool, DB table: talent_pool
// Dev must implement: RegisterRoutes(rg *gin.RouterGroup, h *Handler)
package talent_pool

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TalentPoolEntry maps to the talent_pool table
type TalentPoolEntry struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ApplicationID uuid.UUID  `json:"application_id" db:"application_id"`
	CandidateID   uuid.UUID  `json:"candidate_id" db:"candidate_id"`
	AddedByID     uuid.UUID  `json:"added_by_id" db:"added_by_id"`
	Note          *string    `json:"note,omitempty" db:"note"`
	Tags          []string   `json:"tags" db:"tags"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// TalentPoolService defines business operations for the talent pool
type TalentPoolService interface {
	// Dev must implement: list all entries (HR/Admin)
	ListEntries(ctx *gin.Context, filter TalentPoolFilter) ([]TalentPoolEntry, int, error)
	// Dev must implement: get a single entry by ID
	GetEntry(ctx *gin.Context, id uuid.UUID) (*TalentPoolEntry, error)
	// Dev must implement: add candidate to talent pool (from Transfer action)
	AddEntry(ctx *gin.Context, req AddTalentPoolRequest) (*TalentPoolEntry, error)
	// Dev must implement: remove candidate from talent pool
	RemoveEntry(ctx *gin.Context, id uuid.UUID) error
	// Dev must implement: update note/tags on a talent pool entry
	UpdateEntry(ctx *gin.Context, id uuid.UUID, req UpdateTalentPoolRequest) (*TalentPoolEntry, error)
}

// TalentPoolFilter holds query parameters for listing talent pool entries
type TalentPoolFilter struct {
	Search string
	Tags   []string
	Page   int
	Limit  int
}

// AddTalentPoolRequest is the request body for POST /api/v1/talent-pool
type AddTalentPoolRequest struct {
	ApplicationID uuid.UUID  `json:"application_id" binding:"required"`
	Note          *string    `json:"note,omitempty"`
	Tags          []string   `json:"tags,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
}

// UpdateTalentPoolRequest is the request body for PUT /api/v1/talent-pool/:id
type UpdateTalentPoolRequest struct {
	Note      *string    `json:"note,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// RegisterRoutes mounts talent pool endpoints onto the provided RouterGroup
// Auth: Bearer JWT required for all routes
// Roles: HR, Admin
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	tp := rg.Group("/talent-pool")
	{
		tp.GET("", h.ListEntries)        // Role: HR, Admin
		tp.POST("", h.AddEntry)          // Role: HR, Admin
		tp.GET("/:id", h.GetEntry)       // Role: HR, Admin
		tp.PUT("/:id", h.UpdateEntry)    // Role: HR, Admin
		tp.DELETE("/:id", h.RemoveEntry) // Role: HR, Admin
	}
}

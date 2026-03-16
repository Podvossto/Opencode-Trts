// Purpose: Register HTTP routes for the requisitions domain
// Ref: API Contract — POST/GET/PUT /api/v1/requisitions, DB table: requisitions
// Dev must implement: RegisterRoutes(rg *gin.RouterGroup, svc RequisitionService)
package requisitions

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequisitionStatus represents the lifecycle state of a headcount requisition
type RequisitionStatus string

const (
	RequisitionStatusDraft     RequisitionStatus = "draft"
	RequisitionStatusPending   RequisitionStatus = "pending_approval"
	RequisitionStatusApproved  RequisitionStatus = "approved"
	RequisitionStatusRejected  RequisitionStatus = "rejected"
	RequisitionStatusFulfilled RequisitionStatus = "fulfilled"
	RequisitionStatusCancelled RequisitionStatus = "cancelled"
)

// Requisition maps to the requisitions table
type Requisition struct {
	ID            uuid.UUID         `json:"id" db:"id"`
	JobID         *uuid.UUID        `json:"job_id,omitempty" db:"job_id"`
	RequestedByID uuid.UUID         `json:"requested_by_id" db:"requested_by_id"`
	ApprovedByID  *uuid.UUID        `json:"approved_by_id,omitempty" db:"approved_by_id"`
	Department    string            `json:"department" db:"department"`
	Position      string            `json:"position" db:"position"`
	Headcount     int               `json:"headcount" db:"headcount"`
	Justification string            `json:"justification" db:"justification"`
	Status        RequisitionStatus `json:"status" db:"status"`
	RejectionNote *string           `json:"rejection_note,omitempty" db:"rejection_note"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
}

// RequisitionService defines the business operations for requisitions
type RequisitionService interface {
	// Dev must implement: list all requisitions (HR/Admin)
	ListRequisitions(ctx *gin.Context) ([]Requisition, error)
	// Dev must implement: get a single requisition by ID
	GetRequisition(ctx *gin.Context, id uuid.UUID) (*Requisition, error)
	// Dev must implement: create a new requisition (HR submits)
	CreateRequisition(ctx *gin.Context, req CreateRequisitionRequest) (*Requisition, error)
	// Dev must implement: approve or reject a requisition (Supervisor/Admin)
	ReviewRequisition(ctx *gin.Context, id uuid.UUID, req ReviewRequisitionRequest) (*Requisition, error)
	// Dev must implement: cancel a requisition
	CancelRequisition(ctx *gin.Context, id uuid.UUID) error
}

// CreateRequisitionRequest is the request body for POST /api/v1/requisitions
type CreateRequisitionRequest struct {
	Department    string `json:"department" binding:"required"`
	Position      string `json:"position" binding:"required"`
	Headcount     int    `json:"headcount" binding:"required,min=1"`
	Justification string `json:"justification" binding:"required"`
}

// ReviewRequisitionRequest is the request body for PUT /api/v1/requisitions/:id/review
type ReviewRequisitionRequest struct {
	Action        string  `json:"action" binding:"required,oneof=approve reject"`
	RejectionNote *string `json:"rejection_note,omitempty"`
}

// RegisterRoutes mounts requisition endpoints onto the provided RouterGroup
// Auth: Bearer JWT required for all routes
// Roles: HR (create), Supervisor/Admin (review), HR/Admin (list/get)
// Dev must implement: inject middleware.AuthMiddleware() and middleware.RequireRole()
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	r := rg.Group("/requisitions")
	{
		r.GET("", h.ListRequisitions)             // Role: HR, Admin, Supervisor
		r.POST("", h.CreateRequisition)           // Role: HR
		r.GET("/:id", h.GetRequisition)           // Role: HR, Admin, Supervisor
		r.PUT("/:id/review", h.ReviewRequisition) // Role: Supervisor, Admin
		r.DELETE("/:id", h.CancelRequisition)     // Role: HR (own), Admin
	}
}

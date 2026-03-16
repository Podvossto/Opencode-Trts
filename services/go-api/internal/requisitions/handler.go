// Purpose: HTTP handler for requisitions domain
// Ref: API Contract — GET/POST/PUT/DELETE /api/v1/requisitions
// Dev must implement: bind request bodies, call service methods, return standard responses
package requisitions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler holds the RequisitionService dependency
type Handler struct {
	svc RequisitionService
}

// NewHandler constructs a requisitions Handler
func NewHandler(svc RequisitionService) *Handler {
	return &Handler{svc: svc}
}

// ListRequisitions handles GET /api/v1/requisitions
// Auth: Bearer JWT | Role: HR, Admin, Supervisor
// Dev must implement: optional ?status= and ?department= query filters
func (h *Handler) ListRequisitions(c *gin.Context) {
	// TODO: parse query params (status, department, page, limit)
	// TODO: call h.svc.ListRequisitions(c)
	// TODO: return response.Success(c, http.StatusOK, result)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// GetRequisition handles GET /api/v1/requisitions/:id
// Auth: Bearer JWT | Role: HR, Admin, Supervisor
func (h *Handler) GetRequisition(c *gin.Context) {
	// TODO: parse c.Param("id") as uuid.UUID
	// TODO: call h.svc.GetRequisition(c, id)
	// TODO: return 404 if not found
	_ = uuid.UUID{}
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// CreateRequisition handles POST /api/v1/requisitions
// Auth: Bearer JWT | Role: HR
// Body: CreateRequisitionRequest
func (h *Handler) CreateRequisition(c *gin.Context) {
	// TODO: bind CreateRequisitionRequest
	// TODO: extract requesting user ID from JWT claims
	// TODO: call h.svc.CreateRequisition(c, req)
	// TODO: return 201 with created requisition
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// ReviewRequisition handles PUT /api/v1/requisitions/:id/review
// Auth: Bearer JWT | Role: Supervisor, Admin
// Body: ReviewRequisitionRequest { action: "approve"|"reject", rejection_note? }
func (h *Handler) ReviewRequisition(c *gin.Context) {
	// TODO: parse c.Param("id") as uuid.UUID
	// TODO: bind ReviewRequisitionRequest
	// TODO: validate rejection_note is present when action == "reject"
	// TODO: call h.svc.ReviewRequisition(c, id, req)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// CancelRequisition handles DELETE /api/v1/requisitions/:id
// Auth: Bearer JWT | Role: HR (own only), Admin
func (h *Handler) CancelRequisition(c *gin.Context) {
	// TODO: parse c.Param("id") as uuid.UUID
	// TODO: verify ownership (HR can only cancel own requisitions)
	// TODO: call h.svc.CancelRequisition(c, id)
	// TODO: return 204 on success
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

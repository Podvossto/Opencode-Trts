// Purpose: HTTP handlers for the Jobs domain — public portal listing
// Ref: API Contract — GET /api/v1/portal/jobs, GET /api/v1/portal/jobs/:id
// ATS-012: Public job listing + detail endpoints (no auth required)
// Full job CRUD handlers (Sprint 2) will be added later
package jobs

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tigersoft/ats-go-api/pkg/response"
)

// PortalHandler holds portal HTTP handlers
type PortalHandler struct {
	svc *PortalService
}

// NewPortalHandler creates a new portal handler
func NewPortalHandler(svc *PortalService) *PortalHandler {
	return &PortalHandler{svc: svc}
}

// ListPublicJobs handles GET /api/v1/portal/jobs
// Query params: page, limit, search
func (h *PortalHandler) ListPublicJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	result, err := h.svc.ListPublicJobs(c.Request.Context(), page, limit, search)
	if err != nil {
		log.Printf("list public jobs error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, result)
}

// GetPublicJob handles GET /api/v1/portal/jobs/:id
func (h *PortalHandler) GetPublicJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	job, err := h.svc.GetPublicJob(c.Request.Context(), id)
	if err != nil {
		log.Printf("get public job error: %v", err)
		response.InternalServerError(c)
		return
	}
	if job == nil {
		response.NotFound(c)
		return
	}

	response.OK(c, job)
}

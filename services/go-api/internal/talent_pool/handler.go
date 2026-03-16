// Purpose: HTTP handler for the talent pool domain
// Ref: API Contract — GET/POST/PUT/DELETE /api/v1/talent-pool
// Dev must implement: bind request bodies, call service methods, return standard responses
package talent_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler holds the TalentPoolService dependency
type Handler struct {
	svc TalentPoolService
}

// NewHandler constructs a talent pool Handler
func NewHandler(svc TalentPoolService) *Handler {
	return &Handler{svc: svc}
}

// ListEntries handles GET /api/v1/talent-pool
// Auth: Bearer JWT | Role: HR, Admin
// Query: ?search=&tags=&page=&limit=
func (h *Handler) ListEntries(c *gin.Context) {
	// TODO: parse query params into TalentPoolFilter
	// TODO: call h.svc.ListEntries(c, filter)
	// TODO: return response.Paginated(c, entries, total, page, limit)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// GetEntry handles GET /api/v1/talent-pool/:id
// Auth: Bearer JWT | Role: HR, Admin
func (h *Handler) GetEntry(c *gin.Context) {
	// TODO: parse c.Param("id") as uuid.UUID
	// TODO: call h.svc.GetEntry(c, id)
	// TODO: return 404 if not found
	_ = uuid.UUID{}
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// AddEntry handles POST /api/v1/talent-pool
// Auth: Bearer JWT | Role: HR, Admin
// Body: AddTalentPoolRequest
func (h *Handler) AddEntry(c *gin.Context) {
	// TODO: bind AddTalentPoolRequest
	// TODO: extract added_by_id from JWT claims
	// TODO: call h.svc.AddEntry(c, req)
	// TODO: return 201 with created entry
	// TODO: return 409 if application already in pool
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// UpdateEntry handles PUT /api/v1/talent-pool/:id
// Auth: Bearer JWT | Role: HR, Admin
// Body: UpdateTalentPoolRequest
func (h *Handler) UpdateEntry(c *gin.Context) {
	// TODO: parse c.Param("id") as uuid.UUID
	// TODO: bind UpdateTalentPoolRequest
	// TODO: call h.svc.UpdateEntry(c, id, req)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// RemoveEntry handles DELETE /api/v1/talent-pool/:id
// Auth: Bearer JWT | Role: HR, Admin
func (h *Handler) RemoveEntry(c *gin.Context) {
	// TODO: parse c.Param("id") as uuid.UUID
	// TODO: call h.svc.RemoveEntry(c, id)
	// TODO: return 204 on success
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

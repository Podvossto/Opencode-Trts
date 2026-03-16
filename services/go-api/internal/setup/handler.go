// Purpose: HTTP handler for system setup / configuration domain
// Ref: API Contract — GET/POST/PUT/DELETE /api/v1/setup/*
// Dev must implement: bind request bodies, call service methods, return standard responses
package setup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler holds the SetupService dependency
type Handler struct {
	svc SetupService
}

// NewHandler constructs a setup Handler
func NewHandler(svc SetupService) *Handler {
	return &Handler{svc: svc}
}

// ListDepartments handles GET /api/v1/setup/departments
// Auth: Bearer JWT | Role: Admin
func (h *Handler) ListDepartments(c *gin.Context) {
	// TODO: call h.svc.ListDepartments(c)
	// TODO: return response.Success(c, http.StatusOK, result)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// CreateDepartment handles POST /api/v1/setup/departments
// Auth: Bearer JWT | Role: Admin
// Body: UpsertDepartmentRequest
func (h *Handler) CreateDepartment(c *gin.Context) {
	// TODO: bind UpsertDepartmentRequest
	// TODO: call h.svc.CreateDepartment(c, req)
	// TODO: return 201 with created department
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// UpdateDepartment handles PUT /api/v1/setup/departments/:id
// Auth: Bearer JWT | Role: Admin
// Body: UpsertDepartmentRequest
func (h *Handler) UpdateDepartment(c *gin.Context) {
	// TODO: parse c.Param("id") as uuid.UUID
	// TODO: bind UpsertDepartmentRequest
	// TODO: call h.svc.UpdateDepartment(c, id, req)
	_ = uuid.UUID{}
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// DeleteDepartment handles DELETE /api/v1/setup/departments/:id
// Auth: Bearer JWT | Role: Admin
func (h *Handler) DeleteDepartment(c *gin.Context) {
	// TODO: parse c.Param("id") as uuid.UUID
	// TODO: call h.svc.DeleteDepartment(c, id)
	// TODO: return 204 on success, 409 if department has active jobs
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// ListEmailTemplates handles GET /api/v1/setup/email-templates
// Auth: Bearer JWT | Role: Admin
func (h *Handler) ListEmailTemplates(c *gin.Context) {
	// TODO: call h.svc.ListEmailTemplates(c)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// GetEmailTemplate handles GET /api/v1/setup/email-templates/:slug
// Auth: Bearer JWT | Role: Admin
func (h *Handler) GetEmailTemplate(c *gin.Context) {
	// TODO: parse c.Param("slug")
	// TODO: call h.svc.GetEmailTemplate(c, slug)
	// TODO: return 404 if not found
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// UpdateEmailTemplate handles PUT /api/v1/setup/email-templates/:slug
// Auth: Bearer JWT | Role: Admin
// Body: UpdateEmailTemplateRequest
func (h *Handler) UpdateEmailTemplate(c *gin.Context) {
	// TODO: parse c.Param("slug")
	// TODO: bind UpdateEmailTemplateRequest
	// TODO: call h.svc.UpdateEmailTemplate(c, slug, req)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// GetSettings handles GET /api/v1/setup/settings
// Auth: Bearer JWT | Role: Admin
func (h *Handler) GetSettings(c *gin.Context) {
	// TODO: call h.svc.GetSettings(c)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

// UpdateSetting handles PUT /api/v1/setup/settings/:key
// Auth: Bearer JWT | Role: Admin
// Body: { "value": "string" }
func (h *Handler) UpdateSetting(c *gin.Context) {
	// TODO: parse c.Param("key")
	// TODO: bind body { value string }
	// TODO: call h.svc.UpdateSetting(c, key, value)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}

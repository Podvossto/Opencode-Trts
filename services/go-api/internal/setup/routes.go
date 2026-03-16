// Purpose: Register HTTP routes for system setup / configuration domain
// Ref: API Contract — GET/PUT /api/v1/setup/*, DB tables: departments, email_templates, system_settings
// Dev must implement: RegisterRoutes(rg *gin.RouterGroup, h *Handler)
package setup

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Department maps to the departments table
type Department struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Code      string     `json:"code" db:"code"`
	ManagerID *uuid.UUID `json:"manager_id,omitempty" db:"manager_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// EmailTemplate maps to the email_templates table
type EmailTemplate struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Slug      string    `json:"slug" db:"slug"`
	Subject   string    `json:"subject" db:"subject"`
	Body      string    `json:"body" db:"body"`
	Variables []string  `json:"variables" db:"variables"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SystemSetting is a generic key-value configuration row
type SystemSetting struct {
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SetupService defines the business operations for system configuration
type SetupService interface {
	// Departments
	ListDepartments(ctx *gin.Context) ([]Department, error)
	CreateDepartment(ctx *gin.Context, req UpsertDepartmentRequest) (*Department, error)
	UpdateDepartment(ctx *gin.Context, id uuid.UUID, req UpsertDepartmentRequest) (*Department, error)
	DeleteDepartment(ctx *gin.Context, id uuid.UUID) error
	// Email templates
	ListEmailTemplates(ctx *gin.Context) ([]EmailTemplate, error)
	GetEmailTemplate(ctx *gin.Context, slug string) (*EmailTemplate, error)
	UpdateEmailTemplate(ctx *gin.Context, slug string, req UpdateEmailTemplateRequest) (*EmailTemplate, error)
	// System settings
	GetSettings(ctx *gin.Context) ([]SystemSetting, error)
	UpdateSetting(ctx *gin.Context, key string, value string) (*SystemSetting, error)
}

// UpsertDepartmentRequest is the request body for POST/PUT /api/v1/setup/departments
type UpsertDepartmentRequest struct {
	Name      string     `json:"name" binding:"required"`
	Code      string     `json:"code" binding:"required"`
	ManagerID *uuid.UUID `json:"manager_id,omitempty"`
}

// UpdateEmailTemplateRequest is the request body for PUT /api/v1/setup/email-templates/:slug
type UpdateEmailTemplateRequest struct {
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// RegisterRoutes mounts setup/configuration endpoints onto the provided RouterGroup
// Auth: Bearer JWT required; Role: Admin only for all routes
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	s := rg.Group("/setup")
	{
		// Departments
		s.GET("/departments", h.ListDepartments)
		s.POST("/departments", h.CreateDepartment)
		s.PUT("/departments/:id", h.UpdateDepartment)
		s.DELETE("/departments/:id", h.DeleteDepartment)
		// Email templates
		s.GET("/email-templates", h.ListEmailTemplates)
		s.GET("/email-templates/:slug", h.GetEmailTemplate)
		s.PUT("/email-templates/:slug", h.UpdateEmailTemplate)
		// System settings
		s.GET("/settings", h.GetSettings)
		s.PUT("/settings/:key", h.UpdateSetting)
	}
}

// Purpose: Forms domain — application form builder CRUD (stores schema as JSONB)
// Ref: API Contract — GET/POST /api/v1/forms, GET/PUT/DELETE /api/v1/forms/{id}
// DB tables: application_forms
// Dev must implement: FormHandler, FormService, FormRepository
package forms

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
)

// RegisterRoutes mounts all form builder routes.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	// Dev must implement:
	// GET  /forms       (HR/Admin)
	// POST /forms       (HR/Admin)
	// GET  /forms/:id   (HR/Admin)
	// PUT  /forms/:id   (HR/Admin)
	// DELETE /forms/:id (HR/Admin)
}

// FieldType is one of 14 supported field types (SRS §Form Builder)
type FieldType string

const (
	FieldText        FieldType = "text"
	FieldTextarea    FieldType = "textarea"
	FieldDropdown    FieldType = "dropdown"
	FieldCheckbox    FieldType = "checkbox"
	FieldRadio       FieldType = "radio"
	FieldFileUpload  FieldType = "file_upload"
	FieldDatePicker  FieldType = "date_picker"
	FieldMCQ         FieldType = "mcq"
	FieldShortAnswer FieldType = "short_answer"
	FieldEssay       FieldType = "essay"
	FieldTrueFalse   FieldType = "true_false"
	FieldRatingScale FieldType = "rating_scale"
	FieldImage       FieldType = "image"
	FieldVideo       FieldType = "video"
	FieldLink        FieldType = "link"
)

// FormField represents a single field in the form schema JSONB
type FormField struct {
	ID       string    `json:"id"`
	Type     FieldType `json:"type"`
	Label    string    `json:"label"`
	Required bool      `json:"required"`
	Options  []string  `json:"options,omitempty"` // for dropdown, checkbox, radio, mcq
	Order    int       `json:"order"`
}

// ApplicationForm is the domain model for a form template
type ApplicationForm struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	FormSchema json.RawMessage `json:"form_schema"` // JSONB: []FormField
	CreatedAt  time.Time       `json:"created_at"`
}

// FormService defines the form builder business logic contract.
type FormService interface {
	List(ctx context.Context) ([]ApplicationForm, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ApplicationForm, error)
	Create(ctx context.Context, req CreateFormRequest) (*ApplicationForm, error)
	Update(ctx context.Context, id uuid.UUID, req CreateFormRequest) (*ApplicationForm, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// CreateFormRequest is the request body for POST /api/v1/forms
type CreateFormRequest struct {
	Name       string      `json:"name" binding:"required"`
	FormSchema []FormField `json:"form_schema" binding:"required,min=1"`
}

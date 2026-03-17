// Purpose: Forms domain — application form builder CRUD (stores schema as JSONB)
// Ref: ATS-030 — Form CRUD with reference check on delete
// DB tables: application_forms
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
	"github.com/tigersoft/ats-go-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	admin := rg.Group("/forms")
	admin.Use(middleware.AuthMiddleware(cfg, rdb))
	admin.Use(middleware.RequireRole("admin", "hr"))
	{
		admin.GET("", h.ListForms)
		admin.POST("", h.CreateForm)
		admin.GET("/:id", h.GetForm)
		admin.PUT("/:id", h.UpdateForm)
		admin.DELETE("/:id", h.DeleteForm)
		admin.POST("/:id/publish", h.PublishForm)
	}
}

// NewRepository creates a new form repository (alias for DI wiring)
func NewRepository(db *pgxpool.Pool) *formRepository {
	return newFormRepository(db)
}

// NewService creates a new form service from a repository
func NewService(repo *formRepository) FormService {
	return &formService{repo: repo}
}

// NewHandler creates a new form HTTP handler
func NewHandler(svc FormService) *FormHandler {
	return NewFormHandler(svc)
}

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

type FormField struct {
	ID       string    `json:"id"`
	Type     FieldType `json:"type"`
	Label    string    `json:"label"`
	Required bool      `json:"required"`
	Options  []string  `json:"options,omitempty"`
	Order    int       `json:"order"`
}

type ApplicationForm struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	FormSchema  json.RawMessage `json:"form_schema"`
	IsPublished bool            `json:"is_published"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type FormService interface {
	List(ctx context.Context, page, limit int, published *bool) ([]ApplicationForm, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ApplicationForm, error)
	Create(ctx context.Context, req CreateFormRequest) (*ApplicationForm, error)
	Update(ctx context.Context, id uuid.UUID, req CreateFormRequest) (*ApplicationForm, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID) (*ApplicationForm, error)
}

type CreateFormRequest struct {
	Name       string      `json:"name" binding:"required"`
	FormSchema []FormField `json:"form_schema" binding:"required,min=1"`
}

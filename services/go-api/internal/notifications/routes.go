// Purpose: Routes, models, and service interface for Notifications domain — email templates + send log
// Ref: API Contract — GET /api/v1/notifications/templates, POST /api/v1/notifications/send
// DB tables: (email_templates, notification_logs — Batch 2 additive migration)
// Dev must implement: NotificationHandler, NotificationService
package notifications

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
)

// RegisterRoutes mounts notification routes.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	// Dev must implement:
	// GET  /notifications/templates      (HR/Admin)
	// POST /notifications/send           (HR — manual trigger)
	// GET  /notifications/logs           (Admin)
}

// EmailTemplate represents an editable notification template.
type EmailTemplate struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Stage     string    `json:"stage"` // pipeline stage trigger
	Subject   string    `json:"subject"`
	BodyHTML  string    `json:"body_html"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationService defines the notification business logic contract.
type NotificationService interface {
	ListTemplates(ctx context.Context) ([]EmailTemplate, error)
	SendToApplicant(ctx context.Context, applicationID uuid.UUID, templateID uuid.UUID, actorID uuid.UUID) error
	GetLogs(ctx context.Context, filters map[string]string) ([]interface{}, error)
}

// Purpose: Pipeline domain — HR action dispatch (move stage, email trigger, hire/drop/transfer)
// Ref: API Contract — PUT /api/v1/applications/{id}/status, POST /api/v1/applications/{id}/email
// DB tables: applications
// Dev must implement: PipelineHandler, PipelineService
package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
	"github.com/tigersoft/ats-go-api/internal/email"
	"github.com/tigersoft/ats-go-api/internal/middleware"
)

// RegisterRoutes mounts pipeline action routes.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	repo := NewRepository(db)
	emailSender := email.NewSMTPSender(email.Config{
		Host: cfg.SMTPHost,
		Port: cfg.SMTPPort,
		User: cfg.SMTPUser,
		Pass: cfg.SMTPPass,
		From: cfg.SMTPFrom,
	})
	svc := NewPipelineService(repo, emailSender)
	h := NewHandler(svc)

	pipeline := rg.Group("/pipeline")
	pipeline.Use(middleware.AuthMiddleware(cfg, rdb))
	pipeline.Use(middleware.RequireRole("admin", "hr"))
	{
		pipeline.POST("/jobs/:jobId/email", h.SendApplicantEmail)
	}

	apps := rg.Group("/applications")
	apps.Use(middleware.AuthMiddleware(cfg, rdb))
	apps.Use(middleware.RequireRole("admin", "hr"))
	{
		apps.PUT("/:id/status", h.UpdateApplicationStatus)
	}
}

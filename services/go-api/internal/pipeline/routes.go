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
)

// RegisterRoutes mounts pipeline action routes.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	// Dev must implement:
	// PUT  /applications/:id/status  — HR moves applicant to next stage
	// POST /applications/:id/email   — HR triggers email to applicant
	// PUT  /applications/:id/transfer — HR transfers applicant to different job
}

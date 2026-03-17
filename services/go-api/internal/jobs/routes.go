// Purpose: Jobs domain — route registration
// Ref: API Contract — GET /api/v1/portal/jobs, GET /api/v1/portal/jobs/:id
// ATS-012: Public portal endpoints (no auth)
// Full job CRUD routes (Sprint 2) will be added later
package jobs

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
	"github.com/tigersoft/ats-go-api/internal/middleware"
)

// RegisterRoutes mounts all job routes onto the given router group.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	repo := NewRepository(db)

	// Portal (public, no auth) — ATS-012
	portalSvc := NewPortalService(repo)
	portalH := NewPortalHandler(portalSvc)

	portal := rg.Group("/portal/jobs")
	{
		portal.GET("", portalH.ListPublicJobs)
		portal.GET("/:id", portalH.GetPublicJob)
	}

	// Protected job routes (JWT auth required)
	jobs := rg.Group("/jobs")
	jobs.Use(middleware.AuthMiddleware(cfg, rdb))
	{
		jobs.GET("/:id/dashboard", portalH.GetDashboardStats)
	}

	// TODO Sprint 2: Wire full job CRUD routes (admin/HR only) here
	// admin := rg.Group("/jobs")
	// admin.Use(middleware.AuthMiddleware(cfg, rdb))
	// admin.Use(middleware.RequireRole("admin", "hr"))
	// { ... }
}

// Purpose: Auth domain — route registration
// Ref: API Contract — /api/v1/auth/*
package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
	"github.com/tigersoft/ats-go-api/internal/middleware"
)

// RegisterRoutes mounts all /api/v1/auth routes onto the given router group.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *goredis.Client, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo, rdb, cfg)
	h := NewHandler(svc)

	auth := rg.Group("/auth")
	{
		// Public endpoints
		auth.POST("/login", h.Login)
		auth.POST("/otp/request", h.RequestOTP)
		auth.POST("/password/reset", h.ResetPasswordOTP)

		// Protected endpoints (require valid JWT)
		protected := auth.Group("")
		protected.Use(middleware.AuthMiddleware(cfg, rdb))
		{
			protected.POST("/logout", h.Logout)
			protected.PUT("/password/change", h.ChangePassword)
		}
	}
}

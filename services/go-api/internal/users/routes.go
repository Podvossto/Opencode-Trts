// Purpose: Users domain — route registration
// Ref: API Contract — /api/v1/admin/users/*, /api/v1/users/me, /api/v1/admin/departments
// All admin routes require AuthMiddleware + RequireRole("admin")
// /users/me routes require AuthMiddleware only (any authenticated user)
package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
	"github.com/tigersoft/ats-go-api/internal/middleware"
)

// RegisterRoutes mounts all user management routes onto the given router group.
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *goredis.Client, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	// Self-service endpoints (any authenticated user)
	me := rg.Group("/users")
	me.Use(middleware.AuthMiddleware(cfg, rdb))
	{
		me.GET("/me", h.GetMe)
		me.PUT("/me", h.UpdateMe)
	}

	// Admin-only endpoints
	admin := rg.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg, rdb))
	admin.Use(middleware.RequireRole("admin"))
	{
		// User CRUD
		admin.GET("/users", h.ListUsers)
		admin.POST("/users", h.CreateUser)
		admin.GET("/users/:id", h.GetUser)
		admin.PUT("/users/:id", h.UpdateUser)
		admin.DELETE("/users/:id", h.DeleteUser)
		admin.POST("/users/:id/reset-password", h.AdminResetPassword)

		// Department CRUD
		admin.GET("/departments", h.ListDepartments)
		admin.POST("/departments", h.CreateDepartment)
	}
}

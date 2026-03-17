// Purpose: Go API entry point — bootstraps config, DB, Redis, HTTP server, and background workers
// Dev must implement: DB connection pool, Redis client, Gin router registration, scheduler goroutine start
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/tigersoft/ats-go-api/config"
	"github.com/tigersoft/ats-go-api/internal/applications"
	"github.com/tigersoft/ats-go-api/internal/auth"
	"github.com/tigersoft/ats-go-api/internal/forms"
	"github.com/tigersoft/ats-go-api/internal/jobs"
	"github.com/tigersoft/ats-go-api/internal/pipeline"
	"github.com/tigersoft/ats-go-api/internal/scheduler"
	"github.com/tigersoft/ats-go-api/internal/users"
	"github.com/tigersoft/ats-go-api/pkg/database"
	"github.com/tigersoft/ats-go-api/pkg/redis"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	rdb, err := redis.Connect(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer rdb.Close()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS — allow configured origins (defaults: localhost:3000 + localhost:3001)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register route groups
	v1 := r.Group("/api/v1")
	auth.RegisterRoutes(v1, db, rdb, cfg)
	users.RegisterRoutes(v1, db, rdb, cfg)
	jobs.RegisterRoutes(v1, db, rdb, cfg)
	forms.RegisterRoutes(v1, db, rdb, cfg)
	applications.RegisterRoutes(v1, db, rdb, cfg)
	pipeline.RegisterRoutes(v1, db, rdb, cfg)

	// Background scheduler — auto-publish/close jobs (1-minute tick)
	go scheduler.Start(db, cfg)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(":" + cfg.Port); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server...")
}

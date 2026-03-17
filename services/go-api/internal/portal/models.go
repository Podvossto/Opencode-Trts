// Purpose: Portal test session domain — re-openable test sessions with heartbeat/auto-save
// Ref: ATS-025 — GET /api/portal/test/:token, POST /api/portal/test/:token/heartbeat, POST /api/portal/test/:token/submit
// DB tables: test_sessions, test_rounds, test_assignments
package portal

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	test := rg.Group("/portal/test")
	{
		test.GET("/:token", h.GetSession)
		test.POST("/:token/heartbeat", h.Heartbeat)
		test.POST("/:token/submit", h.Submit)
	}
}

type TestSession struct {
	ID              uuid.UUID       `json:"id"`
	Token           string          `json:"token"`
	Status          string          `json:"status"`
	Answers         json.RawMessage `json:"answers,omitempty"`
	LastHeartbeatAt *time.Time      `json:"last_heartbeat_at,omitempty"`
	SubmittedAt     *time.Time      `json:"submitted_at,omitempty"`
	TestRoundID     uuid.UUID       `json:"test_round_id"`
	Questions       []TestQuestion  `json:"questions,omitempty"`
}

type TestQuestion struct {
	ID      string          `json:"id"`
	Content string          `json:"content"`
	Type    string          `json:"type"`
	Options json.RawMessage `json:"options,omitempty"`
}

type SessionResponse struct {
	SessionID    uuid.UUID       `json:"session_id"`
	Questions    []TestQuestion  `json:"questions"`
	SavedAnswers json.RawMessage `json:"saved_answers"`
}

type HeartbeatRequest struct {
	Answers map[string]interface{} `json:"answers,omitempty"`
}

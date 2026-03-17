// Purpose: Portal test session HTTP handlers
package portal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tigersoft/ats-go-api/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetSession(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.BadRequest(c, []response.ValidationError{{Field: "token", Message: "token is required"}})
		return
	}

	result, status, err := h.svc.GetSession(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			response.NotFound(c)
			return
		}
		log.Printf("get session error: %v", err)
		response.InternalServerError(c)
		return
	}

	if status == "submitted" || status == "expired" {
		c.JSON(http.StatusGone, response.ErrorResponse{Error: "session " + status})
		return
	}

	response.OK(c, SessionResponse{
		SessionID:    result.SessionID,
		Questions:    result.Questions,
		SavedAnswers: result.SavedAnswers,
	})
}

func (h *Handler) Heartbeat(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.BadRequest(c, []response.ValidationError{{Field: "token", Message: "token is required"}})
		return
	}

	var req HeartbeatRequest
	var answers json.RawMessage
	if c.ShouldBindJSON(&req) == nil && req.Answers != nil {
		var err error
		answers, err = json.Marshal(req.Answers)
		if err != nil {
			log.Printf("marshal answers error: %v", err)
			response.InternalServerError(c)
			return
		}
	}

	err := h.svc.Heartbeat(c.Request.Context(), token, answers)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			response.NotFound(c)
			return
		}
		if errors.Is(err, ErrSessionSubmitted) {
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: "session already submitted"})
			return
		}
		log.Printf("heartbeat error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, gin.H{"message": "heartbeat received"})
}

func (h *Handler) Submit(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.BadRequest(c, []response.ValidationError{{Field: "token", Message: "token is required"}})
		return
	}

	var req HeartbeatRequest
	var answers json.RawMessage
	if c.ShouldBindJSON(&req) == nil && req.Answers != nil {
		var err error
		answers, err = json.Marshal(req.Answers)
		if err != nil {
			log.Printf("marshal answers error: %v", err)
			response.InternalServerError(c)
			return
		}
	}

	err := h.svc.Submit(c.Request.Context(), token, answers)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			response.NotFound(c)
			return
		}
		if errors.Is(err, ErrSessionSubmitted) {
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: "session already submitted"})
			return
		}
		log.Printf("submit error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, gin.H{"message": "test submitted"})
}

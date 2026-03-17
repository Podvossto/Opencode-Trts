// Purpose: HTTP handlers for Pipeline actions — stage transitions, email triggers, hire/drop/transfer
// Ref: API Contract — PUT /api/v1/applications/{id}/status, POST /api/v1/applications/{id}/email
//
//	PUT /api/v1/applications/{id}/transfer
//
// DB tables: applications
// Dev must implement: UpdateApplicationStatus, SendApplicantEmail, TransferApplication
package pipeline

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc PipelineService
}

func NewHandler(svc PipelineService) *Handler {
	return &Handler{svc: svc}
}

type SendApplicantEmailRequest struct {
	ApplicantIDs  []string `json:"applicant_ids" binding:"required"`
	TemplateCode  string   `json:"template_code" binding:"required"`
	CustomSubject *string  `json:"custom_subject,omitempty"`
	CustomBody    *string  `json:"custom_body,omitempty"`
}

type UpdateApplicationStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Stage  string `json:"stage" binding:"required"`
	Round  int    `json:"round"`
}

type SendEmailResponse struct {
	Sent    int               `json:"sent"`
	Failed  int               `json:"failed"`
	Results []EmailSendResult `json:"results"`
}

type UpdateStatusResponse struct {
	ApplicationID string `json:"application_id"`
	Status        string `json:"status"`
	Stage         string `json:"stage"`
}

// SendApplicantEmail handles POST /jobs/:jobId/pipeline/email
func (h *Handler) SendApplicantEmail(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jobId is required"})
		return
	}

	var req SendApplicantEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.SendBulkEmail(c.Request.Context(), jobID, SendEmailRequest{
		ApplicantIDs:  req.ApplicantIDs,
		TemplateCode:  req.TemplateCode,
		CustomSubject: req.CustomSubject,
		CustomBody:    req.CustomBody,
	})

	if err != nil {
		if errors.Is(err, ErrTemplateNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SendEmailResponse{
		Sent:    result.Sent,
		Failed:  result.Failed,
		Results: result.Results,
	})
}

// UpdateApplicationStatus handles PUT /applications/:id/status
func (h *Handler) UpdateApplicationStatus(c *gin.Context) {
	appID := c.Param("id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "application id is required"})
		return
	}

	_, err := uuid.Parse(appID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}

	var req UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.UpdateApplicationStatus(c.Request.Context(), appID, req.Status, req.Stage, req.Round)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UpdateStatusResponse{
		ApplicationID: appID,
		Status:        req.Status,
		Stage:         req.Stage,
	})
}

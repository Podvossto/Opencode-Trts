// Purpose: HTTP handlers for the Applications domain — portal submission, autosave, status view
// Ref: API Contract — POST /api/v1/portal/applications, GET /api/v1/portal/applications/{id}
//
//	POST /api/v1/portal/applications/autosave, GET /api/v1/applications
//
// DB tables: candidates, applications, application_documents, test_autosave
// Dev must implement: SubmitApplication, GetApplicationStatus, AutosaveApplication, ListApplications, GetApplication
// IMPORTANT: SubmitApplication returns 201 immediately — OCR is enqueued async via ocr_worker
package applications

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tigersoft/ats-go-api/config"
)

type ApplicationHandler struct {
	svc ApplicationService
	cfg *config.Config
}

func NewApplicationHandler(svc ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		svc: svc,
	}
}

func (h *ApplicationHandler) SubmitApplication(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}

	jobIDStr := c.PostForm("job_id")
	if jobIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_id is required"})
		return
	}
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id format"})
		return
	}

	email := c.PostForm("email")
	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	phone := c.PostForm("phone")
	formDataStr := c.PostForm("form_data")

	if email == "" || firstName == "" || lastName == "" || phone == "" || formDataStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	var formData map[string]interface{}
	if err := parseJSON(formDataStr, &formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form_data format"})
		return
	}

	req := SubmitApplicationRequest{
		JobID:     jobID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		FormData:  formData,
	}

	file, header, err := c.Request.FormFile("resume")
	if err != nil {
		app, err := h.svc.Submit(c.Request.Context(), req)
		if err != nil {
			h.handleSubmitError(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"application_id": app.ID,
			"candidate_id":   app.CandidateID,
			"message":        "Application submitted successfully",
		})
		return
	}
	defer file.Close()

	fileExt := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}
	filePath := filepath.Join(uploadDir, newFileName)
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	fileSizeKB := int32(header.Size / 1024)
	if header.Size%1024 > 0 {
		fileSizeKB++
	}

	app, err := h.svc.SubmitWithDocument(c.Request.Context(), req, "resume", filePath, header.Filename, &fileSizeKB)
	if err != nil {
		h.handleSubmitError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"application_id": app.ID,
		"candidate_id":   app.CandidateID,
		"message":        "Application submitted successfully",
	})
}

func (h *ApplicationHandler) handleSubmitError(c *gin.Context, err error) {
	if err == ErrBlacklisted || contains(err.Error(), "blacklisted") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Application blocked"})
		return
	}
	if err == ErrDuplicateApply || contains(err.Error(), "already applied") {
		c.JSON(http.StatusConflict, gin.H{"error": "Already applied to this job"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (h *ApplicationHandler) GetApplicationStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}

	status, err := h.svc.GetStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"application_id": status.ID,
		"status":         status.Status,
		"stage":          status.Stage,
	})
}

type AutosaveRequest struct {
	Answers map[string]interface{} `json:"answers" binding:"required"`
}

func (h *ApplicationHandler) AutosaveApplication(c *gin.Context) {
	idStr := c.Param("id")
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	jobID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	var req AutosaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.svc.Autosave(c.Request.Context(), jobID, email, req.Answers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to autosave"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"saved_at": time.Now().UTC().Format(time.RFC3339),
	})
}

type ListApplicationsQuery struct {
	JobID  string `form:"job_id"`
	Status string `form:"status"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

func (h *ApplicationHandler) ListApplications(c *gin.Context) {
	var query ListApplicationsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query params"})
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}

	var jobID *uuid.UUID
	if query.JobID != "" {
		id, err := uuid.Parse(query.JobID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id"})
			return
		}
		jobID = &id
	}

	var statusFilter *string
	if query.Status != "" {
		statusFilter = &query.Status
	}

	apps, total, err := h.svc.ListForPipeline(c.Request.Context(), jobID, statusFilter, query.Page, query.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  apps,
		"total": total,
		"page":  query.Page,
		"limit": query.Limit,
	})
}

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}

	app, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil || app == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	c.JSON(http.StatusOK, app)
}

func parseJSON(s string, v interface{}) error {
	if s == "" {
		return nil
	}
	return jsonParse(s, v)
}

func jsonParse(s string, v interface{}) error {
	imported := false
	for _, r := range []rune(s) {
		if r == '{' {
			imported = true
			break
		}
	}
	if !imported {
		return fmt.Errorf("invalid json")
	}
	return parseJSONSimple(s, v)
}

func parseJSONSimple(s string, v interface{}) error {
	m, ok := v.(*map[string]interface{})
	if !ok {
		return fmt.Errorf("cannot parse into non-map type")
	}
	*m = make(map[string]interface{})
	parts := extractJSONParts(s)
	for _, part := range parts {
		if len(part) < 2 {
			continue
		}
		keyVal := splitAtFirst(part, ":")
		if len(keyVal) == 2 {
			key := trimQuotes(keyVal[0])
			val := trimQuotes(keyVal[1])
			(*m)[key] = val
		}
	}
	return nil
}

func extractJSONParts(s string) []string {
	var parts []string
	var current []rune
	inString := false
	depth := 0

	for _, r := range s {
		if r == '"' && (len(current) == 0 || current[len(current)-1] != '\\') {
			inString = !inString
		}
		if !inString {
			if r == '{' || r == '[' {
				depth++
			}
			if r == '}' || r == ']' {
				depth--
			}
			if r == ',' && depth == 0 {
				parts = append(parts, string(current))
				current = nil
				continue
			}
		}
		current = append(current, r)
	}
	if len(current) > 0 {
		parts = append(parts, string(current))
	}
	return parts
}

func splitAtFirst(s, sep string) []string {
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			return []string{s[:i], s[i+len(sep):]}
		}
	}
	return []string{s}
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

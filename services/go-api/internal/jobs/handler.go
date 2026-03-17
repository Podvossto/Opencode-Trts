// Purpose: HTTP handlers for the Jobs domain — public portal listing
// Ref: API Contract — GET /api/v1/portal/jobs, GET /api/v1/portal/jobs/:id
// ATS-012: Public job listing + detail endpoints (no auth required)
// ATS-026: Job CRUD handlers
package jobs

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tigersoft/ats-go-api/pkg/response"
)

// PortalHandler holds portal HTTP handlers
type PortalHandler struct {
	svc *PortalService
}

// NewPortalHandler creates a new portal handler
func NewPortalHandler(svc *PortalService) *PortalHandler {
	return &PortalHandler{svc: svc}
}

// ListPublicJobs handles GET /api/v1/portal/jobs
// Query params: page, limit, search
func (h *PortalHandler) ListPublicJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	result, err := h.svc.ListPublicJobs(c.Request.Context(), page, limit, search)
	if err != nil {
		log.Printf("list public jobs error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, result)
}

// GetPublicJob handles GET /api/v1/portal/jobs/:id
func (h *PortalHandler) GetPublicJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	job, err := h.svc.GetPublicJob(c.Request.Context(), id)
	if err != nil {
		log.Printf("get public job error: %v", err)
		response.InternalServerError(c)
		return
	}
	if job == nil {
		response.NotFound(c)
		return
	}

	response.OK(c, job)
}

// GetDashboardStats handles GET /api/v1/jobs/:id/dashboard (ATS-028)
func (h *PortalHandler) GetDashboardStats(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	stats, err := h.svc.GetDashboardStats(c.Request.Context(), id)
	if err != nil {
		log.Printf("get dashboard stats error: %v", err)
		response.InternalServerError(c)
		return
	}
	if stats == nil {
		response.NotFound(c)
		return
	}

	response.OK(c, stats)
}

// JobHandler holds job CRUD HTTP handlers
type JobHandler struct {
	svc *JobService
}

// NewJobHandler creates a new job handler
func NewJobHandler(svc *JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

// getUserID extracts user ID from context (set by auth middleware)
func getUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// CreateJob handles POST /api/v1/jobs
func (h *JobHandler) CreateJob(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}

	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	skills := make([]JobHardSkill, len(req.HardSkills))
	for i, hs := range req.HardSkills {
		skills[i] = JobHardSkill{
			SkillName: hs.SkillName,
			Weight:    hs.Weight,
		}
	}

	job := Job{
		Title:             req.Title,
		DepartmentID:      req.DepartmentID,
		EmploymentTypeID:  req.EmploymentTypeID,
		ExperienceLevelID: req.ExperienceLevelID,
		FormID:            req.FormID,
		Description:       req.Description,
		Requirements:      req.Requirements,
		Slots:             req.Slots,
		Status:            "draft",
		PublishAt:         req.PublishAt,
		CloseAt:           req.CloseAt,
		CreatedBy:         userID,
	}

	createdJob, err := h.svc.CreateJob(c.Request.Context(), job, skills)
	if err != nil {
		log.Printf("create job error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.Created(c, createdJob)
}

// ListJobs handles GET /api/v1/jobs
func (h *JobHandler) ListJobs(c *gin.Context) {
	filter := JobFilter{
		Status:     c.Query("status"),
		Department: c.Query("department"),
		Search:     c.Query("search"),
		Page:       1,
		Limit:      10,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filter.Page = page
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil && limit > 0 {
		filter.Limit = limit
	}

	result, err := h.svc.ListJobs(c.Request.Context(), filter)
	if err != nil {
		log.Printf("list jobs error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, result)
}

// GetJobByID handles GET /api/v1/jobs/:id
func (h *JobHandler) GetJobByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	job, err := h.svc.GetJobByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("get job error: %v", err)
		response.InternalServerError(c)
		return
	}
	if job == nil {
		response.NotFound(c)
		return
	}

	response.OK(c, job)
}

// UpdateJob handles PUT /api/v1/jobs/:id
func (h *JobHandler) UpdateJob(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "body", Message: err.Error()}})
		return
	}

	skills := make([]JobHardSkill, len(req.HardSkills))
	for i, hs := range req.HardSkills {
		skills[i] = JobHardSkill{
			SkillName: hs.SkillName,
			Weight:    hs.Weight,
		}
	}

	job := Job{
		Title:             req.Title,
		DepartmentID:      req.DepartmentID,
		EmploymentTypeID:  req.EmploymentTypeID,
		ExperienceLevelID: req.ExperienceLevelID,
		FormID:            req.FormID,
		Description:       req.Description,
		Requirements:      req.Requirements,
		Slots:             req.Slots,
		Status:            req.Status,
		PublishAt:         req.PublishAt,
		CloseAt:           req.CloseAt,
	}

	updatedJob, err := h.svc.UpdateJob(c.Request.Context(), id, job, skills)
	if err != nil {
		log.Printf("update job error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, updatedJob)
}

// DeleteJob handles DELETE /api/v1/jobs/:id
func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		response.BadRequest(c, []response.ValidationError{{Field: "id", Message: "invalid UUID"}})
		return
	}

	err := h.svc.DeleteJob(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "cannot delete open job" {
			response.UnprocessableEntity(c, []response.ValidationError{{Field: "status", Message: "cannot delete open job"}})
			return
		}
		log.Printf("delete job error: %v", err)
		response.InternalServerError(c)
		return
	}

	response.NoContent(c)
}

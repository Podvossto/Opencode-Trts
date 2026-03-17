// Purpose: Jobs domain — shared types and request/response models
// DB table: jobs, job_hard_skills
package jobs

import (
	"time"

	"github.com/google/uuid"
)

// Job is the domain model for a job posting
type Job struct {
	ID                uuid.UUID  `json:"id"`
	Title             string     `json:"title"`
	DepartmentID      uuid.UUID  `json:"department_id"`
	EmploymentTypeID  uuid.UUID  `json:"employment_type_id"`
	ExperienceLevelID uuid.UUID  `json:"experience_level_id"`
	FormID            uuid.UUID  `json:"form_id"`
	Description       string     `json:"description"`
	Requirements      string     `json:"requirements"`
	Slots             int        `json:"slots"`
	Status            string     `json:"status"` // draft | open | closed
	PublishAt         *time.Time `json:"publish_at,omitempty"`
	CloseAt           *time.Time `json:"close_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         uuid.UUID  `json:"created_by"`
}

// CreateJobRequest is the request body for POST /api/v1/jobs
type CreateJobRequest struct {
	Title             string            `json:"title" binding:"required"`
	DepartmentID      uuid.UUID         `json:"department_id" binding:"required"`
	EmploymentTypeID  uuid.UUID         `json:"employment_type_id" binding:"required"`
	ExperienceLevelID uuid.UUID         `json:"experience_level_id" binding:"required"`
	FormID            uuid.UUID         `json:"form_id" binding:"required"`
	Description       string            `json:"description" binding:"required"`
	Requirements      string            `json:"requirements"`
	Slots             int               `json:"slots" binding:"required,min=1"`
	Status            string            `json:"status"`
	PublishAt         *time.Time        `json:"publish_at"`
	CloseAt           *time.Time        `json:"close_at"`
	HardSkills        []HardSkillWeight `json:"hard_skills"`
}

// HardSkillWeight pairs a skill name with a 1-10 weight
type HardSkillWeight struct {
	SkillName string `json:"skill_name" binding:"required"`
	Weight    int    `json:"weight" binding:"required,min=1,max=10"`
}

// JobHardSkill represents a hard skill requirement for a job
type JobHardSkill struct {
	ID        uuid.UUID `json:"id"`
	JobID     uuid.UUID `json:"job_id"`
	SkillName string    `json:"skill_name"`
	Weight    int       `json:"weight"`
}

// JobFilter holds filter parameters for listing jobs
type JobFilter struct {
	Status     string
	Department string
	Search     string
	Page       int
	Limit      int
}

// UpdateJobRequest is the request body for PUT /api/v1/jobs/:id
type UpdateJobRequest = CreateJobRequest

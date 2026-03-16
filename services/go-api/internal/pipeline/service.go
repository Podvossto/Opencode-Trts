// Purpose: Business logic for Pipeline actions — status transitions, email dispatch, transfer
// Ref: API Contract — PipelineService (to be defined here)
// DB tables: applications
// Dev must implement: PipelineService interface + pipelineService struct
// Key: action-based flow: HR freely selects status; email is triggered via internal/email package
package pipeline

import (
	"context"

	"github.com/google/uuid"
)

// PipelineService defines the pipeline action business logic contract.
type PipelineService interface {
	UpdateStatus(ctx context.Context, applicationID uuid.UUID, status string, actorID uuid.UUID) error
	SendEmail(ctx context.Context, applicationID uuid.UUID, templateName string, actorID uuid.UUID) error
	Transfer(ctx context.Context, applicationID uuid.UUID, targetJobID uuid.UUID, actorID uuid.UUID) error
}

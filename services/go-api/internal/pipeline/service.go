// Purpose: Business logic for Pipeline actions — status transitions, email dispatch, transfer
// Ref: API Contract — PipelineService (to be defined here)
// DB tables: applications
// Dev must implement: PipelineService interface + pipelineService struct
// Key: action-based flow: HR freely selects status; email is triggered via internal/email package
package pipeline

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	emailService "github.com/tigersoft/ats-go-api/internal/email"
)

var ErrTemplateNotFound = errors.New("email template not found")

// PipelineService defines the pipeline action business logic contract.
type PipelineService interface {
	UpdateStatus(ctx context.Context, applicationID uuid.UUID, status string, stage string, round int) error
	SendBulkEmail(ctx context.Context, jobID string, req SendEmailRequest) (SendEmailResult, error)
	UpdateApplicationStatus(ctx context.Context, appID, status, stage string, round int) error
}

// SendEmailRequest is the request body for sending bulk emails
type SendEmailRequest struct {
	ApplicantIDs  []string `json:"applicant_ids"`
	TemplateCode  string   `json:"template_code"`
	CustomSubject *string  `json:"custom_subject,omitempty"`
	CustomBody    *string  `json:"custom_body,omitempty"`
}

// EmailSendResult represents the result of sending an email to one applicant
type EmailSendResult struct {
	ApplicationID uuid.UUID `json:"application_id"`
	CandidateID   uuid.UUID `json:"candidate_id"`
	Email         string    `json:"email"`
	Success       bool      `json:"success"`
	Error         string    `json:"error,omitempty"`
}

// SendEmailResult is the aggregate result of bulk email sending
type SendEmailResult struct {
	Sent    int               `json:"sent"`
	Failed  int               `json:"failed"`
	Results []EmailSendResult `json:"results"`
}

// pipelineService is the concrete implementation of PipelineService
type pipelineService struct {
	repo        *pipelineRepository
	emailSender emailService.Sender
}

// NewPipelineService creates a new pipeline service
func NewPipelineService(repo *pipelineRepository, emailSender emailService.Sender) PipelineService {
	return &pipelineService{
		repo:        repo,
		emailSender: emailSender,
	}
}

// SendBulkEmail sends emails to multiple applicants
func (s *pipelineService) SendBulkEmail(ctx context.Context, jobID string, req SendEmailRequest) (SendEmailResult, error) {
	template, err := s.repo.GetTemplateByCode(ctx, req.TemplateCode)
	if err != nil {
		return SendEmailResult{}, err
	}
	if template == nil {
		return SendEmailResult{}, ErrTemplateNotFound
	}

	apps, err := s.repo.GetApplicationsWithCandidates(ctx, jobID, req.ApplicantIDs)
	if err != nil {
		return SendEmailResult{}, err
	}

	result := SendEmailResult{
		Results: make([]EmailSendResult, 0, len(apps)),
	}

	actorID := uuid.Nil
	userID := &actorID

	for _, awc := range apps {
		subject := template.Subject
		body := template.Body

		if req.CustomSubject != nil {
			subject = *req.CustomSubject
		}
		if req.CustomBody != nil {
			body = *req.CustomBody
		}

		body = replaceTemplateVars(body, awc.Candidate.FirstName, awc.Candidate.LastName)

		sendResult := EmailSendResult{
			ApplicationID: awc.Application.ID,
			CandidateID:   awc.Candidate.ID,
			Email:         awc.Candidate.Email,
			Success:       true,
		}

		err := s.emailSender.Send(awc.Candidate.Email, subject, body)
		if err != nil {
			sendResult.Success = false
			sendResult.Error = err.Error()
			result.Failed++
		} else {
			result.Sent++
		}

		candidateID := awc.Candidate.ID
		log := NotificationLog{
			ID:          uuid.New(),
			UserID:      userID,
			CandidateID: &candidateID,
			Channel:     "email",
			Message:     "Subject: " + subject,
			SentAt:      time.Now(),
			Status:      "sent",
		}
		if !sendResult.Success {
			log.Status = "failed"
		}
		_ = s.repo.InsertNotificationLog(ctx, log)

		result.Results = append(result.Results, sendResult)
	}

	return result, nil
}

// UpdateStatus updates an application's status
func (s *pipelineService) UpdateStatus(ctx context.Context, applicationID uuid.UUID, status string, stage string, round int) error {
	return s.repo.UpdateApplicationStatus(ctx, applicationID.String(), status, stage, round)
}

// UpdateApplicationStatus updates an application's status (alias for UpdateStatus)
func (s *pipelineService) UpdateApplicationStatus(ctx context.Context, appID, status, stage string, round int) error {
	return s.repo.UpdateApplicationStatus(ctx, appID, status, stage, round)
}

func replaceTemplateVars(body, firstName, lastName string) string {
	result := body
	result = strings.Replace(result, "{{first_name}}", firstName, -1)
	result = strings.Replace(result, "{{last_name}}", lastName, -1)
	return result
}

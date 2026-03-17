// Purpose: Business logic implementation for the Applications domain
// Ref: API Contract — ApplicationService interface defined in routes.go
// DB tables: candidates, applications, application_documents
// Dev must implement: applicationService struct satisfying ApplicationService interface
// Key: blacklist check (applications.status='dropped'), duplicate check (same candidate+job+open)
package applications

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tigersoft/ats-go-api/internal/ocr_worker"
)

var (
	ErrBlacklisted    = errors.New("application blocked")
	ErrDuplicateApply = errors.New("already applied to this job")
)

type applicationService struct {
	repo      RepositoryInterface
	ocrWorker *ocr_worker.Worker
}

func NewApplicationService(repo RepositoryInterface, ocrWorker *ocr_worker.Worker) ApplicationService {
	return &applicationService{
		repo:      repo,
		ocrWorker: ocrWorker,
	}
}

func (s *applicationService) Submit(ctx context.Context, req SubmitApplicationRequest) (*Application, error) {
	candidate, err := s.repo.FindCandidateByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("find candidate: %w", err)
	}

	if candidate != nil && candidate.IsBlacklisted {
		reason := "Candidate is blacklisted"
		if candidate.BlacklistReason != nil && *candidate.BlacklistReason != "" {
			reason = *candidate.BlacklistReason
		}
		return nil, fmt.Errorf("%w: %s", ErrBlacklisted, reason)
	}

	if candidate == nil {
		candidate, err = s.repo.CreateCandidate(ctx, req.Email, req.FirstName, req.LastName, req.Phone)
		if err != nil {
			return nil, fmt.Errorf("create candidate: %w", err)
		}
	}

	existingApp, err := s.repo.GetApplicationByCandidateAndJob(ctx, candidate.ID, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("check duplicate application: %w", err)
	}
	if existingApp != nil {
		return nil, ErrDuplicateApply
	}

	job, err := s.repo.GetJobByID(ctx, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("get job: %w", err)
	}
	if job == nil {
		return nil, errors.New("job not found")
	}

	app, err := s.repo.CreateApplication(ctx, candidate.ID, req.JobID, req.FormData)
	if err != nil {
		return nil, fmt.Errorf("create application: %w", err)
	}

	return &Application{
		ID:          app.ID,
		CandidateID: app.CandidateID,
		JobID:       app.JobID,
		Status:      ApplicationStatus(app.Status),
		FormData:    app.FormData,
	}, nil
}

func (s *applicationService) SubmitWithDocument(ctx context.Context, req SubmitApplicationRequest, fileType, fileURL, fileName string, fileSizeKB *int32) (*Application, error) {
	app, err := s.Submit(ctx, req)
	if err != nil {
		return nil, err
	}

	doc, err := s.repo.CreateApplicationDocument(ctx, app.ID, fileType, fileURL, fileName, fileSizeKB)
	if err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	if s.ocrWorker != nil && fileType == "resume" {
		err = s.repo.UpdateApplicationDocumentOCRStatus(ctx, doc.ID, "processing", nil)
		if err != nil {
			return nil, fmt.Errorf("update ocr status: %w", err)
		}

		go s.ocrWorker.Enqueue(ocr_worker.OCRJob{
			ApplicationID: app.ID,
			DocumentURL:   fileURL,
		})
	}

	return app, nil
}

func (s *applicationService) GetByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	app, err := s.repo.GetApplicationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get application: %w", err)
	}
	if app == nil {
		return nil, nil
	}

	return &Application{
		ID:          app.ID,
		CandidateID: app.CandidateID,
		JobID:       app.JobID,
		Status:      ApplicationStatus(app.Status),
		FormData:    app.FormData,
	}, nil
}

func (s *applicationService) Autosave(ctx context.Context, jobID uuid.UUID, email string, data map[string]interface{}) error {
	return errors.New("not implemented")
}

func (s *applicationService) ListForPipeline(ctx context.Context, jobID uuid.UUID, filters map[string]string) ([]Application, error) {
	return nil, errors.New("not implemented")
}

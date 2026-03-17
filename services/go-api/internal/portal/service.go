// Purpose: Portal test session service — business logic
package portal

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound  = errors.New("session not found")
	ErrSessionExpired   = errors.New("session expired")
	ErrSessionSubmitted = errors.New("session already submitted")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type GetSessionResult struct {
	SessionID    uuid.UUID
	Questions    []TestQuestion
	SavedAnswers json.RawMessage
}

func (s *Service) GetSession(ctx context.Context, token string) (*GetSessionResult, string, error) {
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, "", err
	}
	if session == nil {
		return nil, "", ErrSessionNotFound
	}

	if session.Status == "submitted" || session.Status == "expired" {
		return nil, session.Status, nil
	}

	questions, err := s.repo.GetQuestionsByTestRoundID(ctx, session.TestRoundID)
	if err != nil {
		return nil, "", err
	}

	var qs []TestQuestion
	for _, q := range questions {
		var opts json.RawMessage
		if q.Options != "" {
			opts = json.RawMessage(q.Options)
		}
		qs = append(qs, TestQuestion{
			ID:      q.ID,
			Content: q.Content,
			Type:    q.Type,
			Options: opts,
		})
	}

	return &GetSessionResult{
		SessionID:    session.ID,
		Questions:    qs,
		SavedAnswers: session.Answers,
	}, "", nil
}

func (s *Service) Heartbeat(ctx context.Context, token string, answers json.RawMessage) error {
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrSessionNotFound
	}
	if session.Status == "submitted" {
		return ErrSessionSubmitted
	}

	return s.repo.UpdateHeartbeat(ctx, session.ID, answers)
}

func (s *Service) Submit(ctx context.Context, token string, answers json.RawMessage) error {
	alreadySubmitted, err := s.repo.CheckSessionAlreadySubmitted(ctx, token)
	if err != nil {
		return err
	}
	if alreadySubmitted {
		return ErrSessionSubmitted
	}

	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrSessionNotFound
	}

	return s.repo.SubmitSession(ctx, session.ID, answers)
}

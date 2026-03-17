package pipeline

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailTemplate struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Slug      string    `json:"slug" db:"slug"`
	Subject   string    `json:"subject" db:"subject"`
	Body      string    `json:"body" db:"body"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Candidate struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
}

type Application struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CandidateID  uuid.UUID `json:"candidate_id" db:"candidate_id"`
	JobID        uuid.UUID `json:"job_id" db:"job_id"`
	Status       string    `json:"status" db:"status"`
	Stage        string    `json:"stage" db:"stage"`
	CurrentRound int       `json:"current_round" db:"current_round"`
}

type ApplicationWithCandidate struct {
	Application Application
	Candidate   Candidate
}

type NotificationLog struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      *uuid.UUID `json:"user_id" db:"user_id"`
	CandidateID *uuid.UUID `json:"candidate_id" db:"candidate_id"`
	Channel     string     `json:"channel" db:"channel"`
	Message     string     `json:"message" db:"message"`
	SentAt      time.Time  `json:"sent_at" db:"sent_at"`
	Status      string     `json:"status" db:"status"`
}

type pipelineRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *pipelineRepository {
	return &pipelineRepository{db: db}
}

func (r *pipelineRepository) GetTemplateByCode(ctx context.Context, code string) (*EmailTemplate, error) {
	var tmpl EmailTemplate
	err := r.db.QueryRow(ctx, `
		SELECT id, slug, subject, body, type, created_at, updated_at 
		FROM email_templates WHERE slug = $1
	`, code).Scan(&tmpl.ID, &tmpl.Slug, &tmpl.Subject, &tmpl.Body, &tmpl.Type, &tmpl.CreatedAt, &tmpl.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (r *pipelineRepository) GetApplicationsWithCandidates(ctx context.Context, jobID string, appIDs []string) ([]ApplicationWithCandidate, error) {
	query := `
		SELECT a.id, a.candidate_id, a.job_id, a.status, a.stage, a.current_round,
		       c.id, c.email, c.first_name, c.last_name
		FROM applications a
		JOIN candidates c ON a.candidate_id = c.id
		WHERE a.job_id = $1
	`
	var args []interface{}
	args = append(args, jobID)

	if len(appIDs) > 0 {
		query += ` AND a.id = ANY($2)`
		args = append(args, appIDs)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ApplicationWithCandidate
	for rows.Next() {
		var awc ApplicationWithCandidate
		err := rows.Scan(
			&awc.Application.ID, &awc.Application.CandidateID, &awc.Application.JobID,
			&awc.Application.Status, &awc.Application.Stage, &awc.Application.CurrentRound,
			&awc.Candidate.ID, &awc.Candidate.Email, &awc.Candidate.FirstName, &awc.Candidate.LastName,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, awc)
	}
	return results, nil
}

func (r *pipelineRepository) InsertNotificationLog(ctx context.Context, log NotificationLog) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO notification_logs (id, user_id, candidate_id, channel, message, sent_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, log.ID, log.UserID, log.CandidateID, log.Channel, log.Message, log.SentAt, log.Status)
	return err
}

func (r *pipelineRepository) UpdateApplicationStatus(ctx context.Context, appID, status, stage string, round int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE applications SET status = $1, stage = $2, current_round = $3 WHERE id = $4
	`, status, stage, round, appID)
	return err
}

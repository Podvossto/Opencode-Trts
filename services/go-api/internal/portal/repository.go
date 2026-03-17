// Purpose: Portal test session repository — database queries
// DB tables: test_sessions, test_rounds, test_assignments
package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

type TestSessionRow struct {
	ID              uuid.UUID
	Token           string
	Status          string
	Answers         json.RawMessage
	LastHeartbeatAt *time.Time
	SubmittedAt     *time.Time
	TestRoundID     uuid.UUID
}

type TestQuestionRow struct {
	ID      string
	Content string
	Type    string
	Options string
}

func (r *Repository) GetSessionByToken(ctx context.Context, token string) (*TestSessionRow, error) {
	var s TestSessionRow
	err := r.db.QueryRow(ctx, `
		SELECT id, token, status, answers, last_heartbeat_at, submitted_at, test_round_id
		FROM test_sessions
		WHERE token = $1`, token).
		Scan(&s.ID, &s.Token, &s.Status, &s.Answers, &s.LastHeartbeatAt, &s.SubmittedAt, &s.TestRoundID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get session by token: %w", err)
	}
	return &s, nil
}

func (r *Repository) GetQuestionsByTestRoundID(ctx context.Context, testRoundID uuid.UUID) ([]TestQuestionRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, content, type, options
		FROM test_questions
		WHERE test_round_id = $1
		ORDER BY order_index`, testRoundID)
	if err != nil {
		return nil, fmt.Errorf("get questions: %w", err)
	}
	defer rows.Close()

	var questions []TestQuestionRow
	for rows.Next() {
		var q TestQuestionRow
		if err := rows.Scan(&q.ID, &q.Content, &q.Type, &q.Options); err != nil {
			return nil, fmt.Errorf("scan question: %w", err)
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *Repository) UpdateHeartbeat(ctx context.Context, sessionID uuid.UUID, answers json.RawMessage) error {
	if answers != nil {
		_, err := r.db.Exec(ctx, `
			UPDATE test_sessions
			SET last_heartbeat_at = NOW(),
			    answers = COALESCE(answers, '{}'::jsonb) || $2::jsonb
			WHERE id = $1`, sessionID, answers)
		if err != nil {
			return fmt.Errorf("update heartbeat: %w", err)
		}
	} else {
		_, err := r.db.Exec(ctx, `
			UPDATE test_sessions
			SET last_heartbeat_at = NOW()
			WHERE id = $1`, sessionID)
		if err != nil {
			return fmt.Errorf("update heartbeat: %w", err)
		}
	}
	return nil
}

func (r *Repository) SubmitSession(ctx context.Context, sessionID uuid.UUID, answers json.RawMessage) error {
	_, err := r.db.Exec(ctx, `
		UPDATE test_sessions
		SET status = 'submitted',
		    submitted_at = NOW(),
		    answers = $2
		WHERE id = $1 AND status != 'submitted'`, sessionID, answers)
	if err != nil {
		return fmt.Errorf("submit session: %w", err)
	}
	return nil
}

func (r *Repository) CheckSessionAlreadySubmitted(ctx context.Context, token string) (bool, error) {
	var status string
	err := r.db.QueryRow(ctx, `SELECT status FROM test_sessions WHERE token = $1`, token).Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("check session status: %w", err)
	}
	return status == "submitted", nil
}

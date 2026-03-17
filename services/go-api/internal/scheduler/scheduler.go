// Purpose: Background scheduler — auto-publish and auto-close job postings (1-minute cron tick)
// Ref: ADR-004 — Go scheduler goroutine
// DB tables: jobs (publish_at, close_at, status)
// IMPORTANT: Auto-close must NOT change status of in-progress applications (SRS FRH10b)
// Dev must implement: Start(db, cfg) — ticker loop, PublishDueJobs(), CloseDueJobs()
package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tigersoft/ats-go-api/config"
)

// Start launches the background scheduler. Call as a goroutine from main.
func Start(ctx context.Context, db *pgxpool.Pool, cfg *config.Config) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	log.Println("[scheduler] started — 60-second tick")

	for {
		select {
		case <-ctx.Done():
			log.Println("[scheduler] stopped")
			return
		case <-ticker.C:
			if err := publishDueJobs(ctx, db); err != nil {
				log.Printf("[scheduler] publishDueJobs error: %v", err)
			}
			if err := closeDueJobs(ctx, db); err != nil {
				log.Printf("[scheduler] closeDueJobs error: %v", err)
			}
		}
	}
}

// publishDueJobs sets status='open' for jobs where publish_at <= now AND status='scheduled'
func publishDueJobs(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		UPDATE jobs 
		SET status = 'open', updated_at = NOW() 
		WHERE status = 'scheduled' AND publish_at <= NOW()`)
	if err != nil {
		return fmt.Errorf("publish due jobs: %w", err)
	}
	return nil
}

// closeDueJobs sets status='closed' for jobs where close_at <= now AND status='open'
// Must NOT affect in-progress applications (applications.status != 'pending')
func closeDueJobs(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		UPDATE jobs 
		SET status = 'closed', updated_at = NOW() 
		WHERE status = 'open' AND close_at <= NOW() AND close_at IS NOT NULL`)
	if err != nil {
		return fmt.Errorf("close due jobs: %w", err)
	}
	return nil
}

// Purpose: Background scheduler — auto-publish and auto-close job postings (1-minute cron tick)
// Ref: ADR-004 — Go scheduler goroutine
// DB tables: jobs (publish_at, close_at, status)
// IMPORTANT: Auto-close must NOT change status of in-progress applications (SRS FRH10b)
// Dev must implement: Start(db, cfg) — ticker loop, PublishDueJobs(), CloseDueJobs()
package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tigersoft/ats-go-api/config"
)

// Start launches the background scheduler. Call as a goroutine from main.
func Start(db *pgxpool.Pool, cfg *config.Config) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("[scheduler] started — 1-minute tick")

	for range ticker.C {
		ctx := context.Background()
		if err := publishDueJobs(ctx, db); err != nil {
			log.Printf("[scheduler] publishDueJobs error: %v", err)
		}
		if err := closeDueJobs(ctx, db); err != nil {
			log.Printf("[scheduler] closeDueJobs error: %v", err)
		}
	}
}

// publishDueJobs sets status='open' for jobs where publish_at <= now AND status='draft'
func publishDueJobs(ctx context.Context, db *pgxpool.Pool) error {
	// Dev must implement: UPDATE jobs SET status='open' WHERE publish_at <= NOW() AND status='draft'
	return nil
}

// closeDueJobs sets status='closed' for jobs where close_at <= now AND status='open'
// Must NOT affect in-progress applications (applications.status != 'pending')
func closeDueJobs(ctx context.Context, db *pgxpool.Pool) error {
	// Dev must implement: UPDATE jobs SET status='closed' WHERE close_at <= NOW() AND status='open'
	return nil
}

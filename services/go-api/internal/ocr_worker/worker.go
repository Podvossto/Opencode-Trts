// Purpose: OCR worker — async goroutine that calls FastAPI /ocr + /embed after application submission
// Ref: ADR-003 — async OCR processing; API Contract — POST /ai/ocr, POST /ai/embed
// Dev must implement: Enqueue(appID, fileURL), worker loop, retry logic (3 attempts)
package ocr_worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OCRJob is the unit of work passed to the worker goroutine
type OCRJob struct {
	ApplicationID uuid.UUID
	DocumentURL   string
}

// Worker processes OCR jobs from a channel
type Worker struct {
	jobs         chan OCRJob
	db           *pgxpool.Pool
	aiServiceURL string
	httpClient   *http.Client
}

// NewWorker creates a new OCR worker. Call Start() as a goroutine.
func NewWorker(db *pgxpool.Pool, aiServiceURL string, bufferSize int) *Worker {
	return &Worker{
		jobs:         make(chan OCRJob, bufferSize),
		db:           db,
		aiServiceURL: aiServiceURL,
		httpClient:   &http.Client{},
	}
}

// Enqueue adds an OCR job to the queue (non-blocking up to bufferSize)
func (w *Worker) Enqueue(job OCRJob) {
	select {
	case w.jobs <- job:
	default:
		log.Printf("[ocr_worker] queue full, dropping job for application %s", job.ApplicationID)
	}
}

// Start begins processing jobs. Call as: go worker.Start(ctx)
func (w *Worker) Start(ctx context.Context) {
	log.Println("[ocr_worker] started")
	for {
		select {
		case <-ctx.Done():
			log.Println("[ocr_worker] stopping")
			return
		case job := <-w.jobs:
			w.process(job)
		}
	}
}

// process calls FastAPI /ocr then /embed; updates application record on success
func (w *Worker) process(job OCRJob) {
	// Dev must implement:
	// 1. POST w.aiServiceURL/ocr with { application_id, document_url }
	// 2. On success, POST w.aiServiceURL/embed with { application_id, text }
	// 3. Update applications SET ocr_status='done' WHERE id = job.ApplicationID
	// 4. Retry up to 3 times on failure, then mark ocr_status='failed'
	log.Printf("[ocr_worker] processing application %s", job.ApplicationID)
	_ = fmt.Sprintf("%s/ocr", w.aiServiceURL)
	_ = bytes.NewBuffer(nil)
	_ = json.Marshal
}

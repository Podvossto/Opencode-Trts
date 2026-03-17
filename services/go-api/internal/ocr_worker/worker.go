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
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxRetries = 3
	ocrTimeout = 30 * time.Second
)

type OCRJob struct {
	ApplicationID uuid.UUID
	DocumentURL   string
}

type Worker struct {
	jobs         chan OCRJob
	db           *pgxpool.Pool
	aiServiceURL string
	httpClient   *http.Client
}

func NewWorker(db *pgxpool.Pool, aiServiceURL string, bufferSize int) *Worker {
	return &Worker{
		jobs:         make(chan OCRJob, bufferSize),
		db:           db,
		aiServiceURL: aiServiceURL,
		httpClient: &http.Client{
			Timeout: ocrTimeout,
		},
	}
}

func (w *Worker) Enqueue(job OCRJob) {
	select {
	case w.jobs <- job:
	default:
		log.Printf("[ocr_worker] queue full, dropping job for application %s", job.ApplicationID)
	}
}

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

func (w *Worker) process(job OCRJob) {
	log.Printf("[ocr_worker] processing application %s", job.ApplicationID)

	var ocrText string
	var success bool

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("[ocr_worker] attempt %d/%d for application %s", attempt, maxRetries, job.ApplicationID)

		text, err := w.callOCR(job.DocumentURL)
		if err != nil {
			log.Printf("[ocr_worker] OCR attempt %d failed for application %s: %v", attempt, job.ApplicationID, err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		ocrText = text
		success = true
		break
	}

	ctx := context.Background()
	if success {
		_, err := w.db.Exec(ctx, `
			UPDATE application_documents
			SET ocr_status = 'done', ocr_text = $1
			WHERE application_id = $2`, ocrText, job.ApplicationID)
		if err != nil {
			log.Printf("[ocr_worker] failed to update OCR status for application %s: %v", job.ApplicationID, err)
		}
		log.Printf("[ocr_worker] OCR completed for application %s", job.ApplicationID)
	} else {
		_, err := w.db.Exec(ctx, `
			UPDATE application_documents
			SET ocr_status = 'failed'
			WHERE application_id = $1`, job.ApplicationID)
		if err != nil {
			log.Printf("[ocr_worker] failed to update OCR status for application %s: %v", job.ApplicationID, err)
		}
		log.Printf("[ocr_worker] OCR failed for application %s after %d attempts", job.ApplicationID, maxRetries)
	}
}

func (w *Worker) callOCR(filePath string) (string, error) {
	reqBody := map[string]string{"file_path": filePath}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", w.aiServiceURL+"/ocr/extract", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call OCR service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OCR service returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	text, ok := result["text"].(string)
	if !ok {
		return "", fmt.Errorf("response missing 'text' field")
	}

	return text, nil
}

func (w *Worker) GetQueueLength() int {
	return len(w.jobs)
}

// Purpose: Server actions for the resume-review feature
// Ref: API Contract — GET /api/v1/resume-review/:jobId, POST /api/v1/resume-review/bulk
// Dev must implement: call Go API via lib/api.ts, handle errors per API contract

'use server';

import type {
  ResumeReviewListResponse,
  BulkReviewRequest,
  ResumeReviewStatus,
} from './types';

/**
 * Fetch all applications pending resume review for a given job
 * GET /api/v1/resume-review/:jobId?page=&limit=
 * Role: HR, Admin
 */
export async function getApplicationsForReview(
  jobId: string,
  page = 1,
  limit = 20
): Promise<ResumeReviewListResponse> {
  // TODO: call api.get(`/resume-review/${jobId}?page=${page}&limit=${limit}`)
  throw new Error('Not implemented');
}

/**
 * Submit a single resume review decision
 * POST /api/v1/resume-review
 * Body: { applicationId, status, notes? }
 * Role: HR, Admin
 */
export async function submitResumeReview(
  applicationId: string,
  status: ResumeReviewStatus,
  notes?: string
): Promise<void> {
  // TODO: call api.post('/resume-review', { applicationId, status, notes })
  throw new Error('Not implemented');
}

/**
 * Bulk update resume review statuses for multiple applications
 * POST /api/v1/resume-review/bulk
 * Body: BulkReviewRequest
 * Role: HR, Admin
 */
export async function bulkSubmitResumeReviews(req: BulkReviewRequest): Promise<void> {
  // TODO: call api.post('/resume-review/bulk', req)
  throw new Error('Not implemented');
}

/**
 * Trigger OCR re-processing for a specific application's resume
 * POST /api/v1/resume-review/:applicationId/ocr
 * Role: HR, Admin
 */
export async function triggerOcrReprocess(applicationId: string): Promise<void> {
  // TODO: call api.post(`/resume-review/${applicationId}/ocr`)
  throw new Error('Not implemented');
}

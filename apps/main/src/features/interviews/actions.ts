// Purpose: Server actions for the interviews feature
// Ref: API Contract — GET/POST/PUT /api/v1/interviews
// Dev must implement: call Go API via lib/api.ts

'use server';

import type {
  InterviewListResponse,
  InterviewSchedule,
  ScheduleInterviewRequest,
  RecordInterviewResultRequest,
} from './types';

/**
 * Fetch all interviews for a specific job
 * GET /api/v1/interviews?jobId=&page=&limit=
 * Role: HR, Admin, Supervisor
 */
export async function listInterviewsByJob(
  jobId: string,
  page = 1,
  limit = 20
): Promise<InterviewListResponse> {
  // TODO: call api.get(`/interviews?jobId=${jobId}&page=${page}&limit=${limit}`)
  throw new Error('Not implemented');
}

/**
 * Schedule a new interview round for an application
 * POST /api/v1/interviews
 * Body: ScheduleInterviewRequest
 * Role: HR, Admin
 */
export async function scheduleInterview(req: ScheduleInterviewRequest): Promise<InterviewSchedule> {
  // TODO: call api.post('/interviews', req)
  throw new Error('Not implemented');
}

/**
 * Record the result of a completed interview
 * PUT /api/v1/interviews/:id/result
 * Body: RecordInterviewResultRequest
 * Role: HR, Admin, Interviewer
 */
export async function recordInterviewResult(
  interviewId: string,
  req: RecordInterviewResultRequest
): Promise<InterviewSchedule> {
  // TODO: call api.put(`/interviews/${interviewId}/result`, req)
  throw new Error('Not implemented');
}

/**
 * Cancel a scheduled interview
 * PUT /api/v1/interviews/:id/cancel
 * Role: HR, Admin
 */
export async function cancelInterview(interviewId: string): Promise<void> {
  // TODO: call api.put(`/interviews/${interviewId}/cancel`)
  throw new Error('Not implemented');
}

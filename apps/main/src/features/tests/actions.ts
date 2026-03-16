// Purpose: Server actions for the tests (online test) feature — HR side
// Ref: API Contract — GET/POST /api/v1/tests
// Dev must implement: call Go API via lib/api.ts

'use server';

import type { TestListResponse, TestAssignment, AssignTestRequest, TestResultDetail } from './types';

/**
 * Fetch all test assignments for a specific job
 * GET /api/v1/tests?jobId=&page=&limit=
 * Role: HR, Admin
 */
export async function listTestsByJob(
  jobId: string,
  page = 1,
  limit = 20
): Promise<TestListResponse> {
  // TODO: call api.get(`/tests?jobId=${jobId}&page=${page}&limit=${limit}`)
  throw new Error('Not implemented');
}

/**
 * Assign an online test to an applicant
 * POST /api/v1/tests
 * Body: AssignTestRequest
 * Role: HR, Admin
 */
export async function assignTest(req: AssignTestRequest): Promise<TestAssignment> {
  // TODO: call api.post('/tests', req)
  // TODO: Go backend will generate a one-time token and send email to candidate
  throw new Error('Not implemented');
}

/**
 * Re-send the test invitation email to the candidate
 * POST /api/v1/tests/:id/resend
 * Role: HR, Admin
 */
export async function resendTestInvitation(assignmentId: string): Promise<void> {
  // TODO: call api.post(`/tests/${assignmentId}/resend`)
  throw new Error('Not implemented');
}

/**
 * Fetch detailed test result with per-question breakdown
 * GET /api/v1/tests/:id/result
 * Role: HR, Admin
 */
export async function getTestResultDetail(assignmentId: string): Promise<TestResultDetail> {
  // TODO: call api.get(`/tests/${assignmentId}/result`)
  throw new Error('Not implemented');
}

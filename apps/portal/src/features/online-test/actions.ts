// Purpose: Server actions for applicant-facing online test delivery (M07)
// Ref: API Contract — GET /api/v1/portal/test/:token, POST /api/v1/portal/test/:token/submit
// Dev must implement: fetchTest, submitAnswers, heartbeat, autoSave

'use server';

import type {
  OnlineTestSession,
  SubmitAnswersRequest,
  AutoSaveRequest,
} from './types';

/**
 * Fetch test session for the applicant (via unique token link)
 * GET /api/v1/portal/test/:token
 * Role: Public (token-authenticated)
 */
export async function fetchTestSession(token: string): Promise<OnlineTestSession> {
  // TODO: call api.get(`/portal/test/${token}`)
  throw new Error('Not implemented');
}

/**
 * Submit test answers
 * POST /api/v1/portal/test/:token/submit
 * Role: Public (token-authenticated)
 */
export async function submitTestAnswers(token: string, data: SubmitAnswersRequest): Promise<void> {
  // TODO: call api.post(`/portal/test/${token}/submit`, data)
  throw new Error('Not implemented');
}

/**
 * Send heartbeat to keep test session alive + report anti-cheat events
 * POST /api/v1/portal/test/:token/heartbeat
 * Role: Public (token-authenticated)
 */
export async function sendHeartbeat(
  token: string,
  data: { tabSwitchCount: number; copyPasteCount: number }
): Promise<void> {
  // TODO: call api.post(`/portal/test/${token}/heartbeat`, data)
  throw new Error('Not implemented');
}

/**
 * Auto-save partial answers (periodic save during test)
 * POST /api/v1/portal/test/:token/autosave
 * Role: Public (token-authenticated)
 */
export async function autoSaveAnswers(token: string, data: AutoSaveRequest): Promise<void> {
  // TODO: call api.post(`/portal/test/${token}/autosave`, data)
  throw new Error('Not implemented');
}

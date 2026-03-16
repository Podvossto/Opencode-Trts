// Purpose: Server actions for the candidates feature
// Ref: API Contract — GET /api/v1/candidates/:id, POST /api/v1/candidates/:id/blacklist
// Dev must implement: call Go API via lib/api.ts, handle 404 and 403 per contract

'use server';

import type { CandidateProfile, CandidateListResponse, BlacklistRequest } from './types';

/**
 * Fetch paginated candidate list
 * GET /api/v1/candidates?search=&page=&limit=
 * Role: HR, Admin
 */
export async function listCandidates(
  search = '',
  page = 1,
  limit = 20
): Promise<CandidateListResponse> {
  // TODO: call api.get(`/candidates?search=${search}&page=${page}&limit=${limit}`)
  throw new Error('Not implemented');
}

/**
 * Fetch the full 360° profile of a single candidate
 * GET /api/v1/candidates/:id
 * Role: HR, Admin
 */
export async function getCandidateProfile(candidateId: string): Promise<CandidateProfile> {
  // TODO: call api.get(`/candidates/${candidateId}`)
  // TODO: handle 404 → redirect to candidates list
  throw new Error('Not implemented');
}

/**
 * Blacklist a candidate (prevents future applications)
 * POST /api/v1/candidates/:id/blacklist
 * Body: BlacklistRequest { reason }
 * Role: HR, Admin
 */
export async function blacklistCandidate(
  candidateId: string,
  req: BlacklistRequest
): Promise<void> {
  // TODO: call api.post(`/candidates/${candidateId}/blacklist`, req)
  throw new Error('Not implemented');
}

/**
 * Remove blacklist status from a candidate
 * DELETE /api/v1/candidates/:id/blacklist
 * Role: Admin
 */
export async function removeBlacklist(candidateId: string): Promise<void> {
  // TODO: call api.delete(`/candidates/${candidateId}/blacklist`)
  throw new Error('Not implemented');
}

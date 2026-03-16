// Purpose: Server actions for the requisitions feature
// Ref: API Contract — GET/POST/PUT/DELETE /api/v1/requisitions
// Dev must implement: call Go API via lib/api.ts

'use server';

import type {
  RequisitionListResponse,
  Requisition,
  CreateRequisitionRequest,
  ReviewRequisitionRequest,
} from './types';

/**
 * Fetch all requisitions (with optional status/department filters)
 * GET /api/v1/requisitions?status=&department=&page=&limit=
 * Role: HR, Admin, Supervisor
 */
export async function listRequisitions(
  status?: string,
  department?: string,
  page = 1,
  limit = 20
): Promise<RequisitionListResponse> {
  // TODO: build query params and call api.get('/requisitions', { params })
  throw new Error('Not implemented');
}

/**
 * Fetch a single requisition by ID
 * GET /api/v1/requisitions/:id
 * Role: HR, Admin, Supervisor
 */
export async function getRequisition(id: string): Promise<Requisition> {
  // TODO: call api.get(`/requisitions/${id}`)
  throw new Error('Not implemented');
}

/**
 * Submit a new headcount requisition
 * POST /api/v1/requisitions
 * Body: CreateRequisitionRequest
 * Role: HR
 */
export async function createRequisition(req: CreateRequisitionRequest): Promise<Requisition> {
  // TODO: call api.post('/requisitions', req)
  throw new Error('Not implemented');
}

/**
 * Approve or reject a pending requisition
 * PUT /api/v1/requisitions/:id/review
 * Body: ReviewRequisitionRequest
 * Role: Supervisor, Admin
 */
export async function reviewRequisition(
  id: string,
  req: ReviewRequisitionRequest
): Promise<Requisition> {
  // TODO: call api.put(`/requisitions/${id}/review`, req)
  throw new Error('Not implemented');
}

/**
 * Cancel a requisition
 * DELETE /api/v1/requisitions/:id
 * Role: HR (own), Admin
 */
export async function cancelRequisition(id: string): Promise<void> {
  // TODO: call api.delete(`/requisitions/${id}`)
  throw new Error('Not implemented');
}

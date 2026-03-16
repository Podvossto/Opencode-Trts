// Purpose: Server actions for the transfer feature
// Ref: API Contract — POST /api/v1/transfer, GET /api/v1/transfer
// Dev must implement: call Go API via lib/api.ts

'use server';

import type { TransferListResponse, TransferRecord, InitiateTransferRequest } from './types';

/**
 * Initiate a transfer for an application to another job or talent pool
 * POST /api/v1/transfer
 * Body: InitiateTransferRequest
 * Role: HR, Admin
 */
export async function initiateTransfer(req: InitiateTransferRequest): Promise<TransferRecord> {
  // TODO: call api.post('/transfer', req)
  // TODO: if destination == 'talent_pool', Go backend also inserts into talent_pool table
  throw new Error('Not implemented');
}

/**
 * Fetch transfer history for a specific application
 * GET /api/v1/transfer?applicationId=
 * Role: HR, Admin
 */
export async function getTransferHistory(applicationId: string): Promise<TransferListResponse> {
  // TODO: call api.get(`/transfer?applicationId=${applicationId}`)
  throw new Error('Not implemented');
}

// Purpose: Portal jobs feature — server actions for public job listing
// Ref: API Contract — GET /api/v1/portal/jobs, GET /api/v1/portal/jobs/:id
// Role: Public (no auth)

import portalApiClient from '@/lib/api'
import type { PublicJobListResponse, PublicJobDetail } from './types'

export interface ListPublicJobsParams {
  page?: number
  limit?: number
  search?: string
}

/**
 * Fetch paginated list of open jobs from the portal API.
 * API response shape: { data: { data: PublicJob[], total, page, limit, total_pages } }
 */
export async function listPublicJobsAction(
  params: ListPublicJobsParams = {}
): Promise<PublicJobListResponse> {
  const { page = 1, limit = 12, search } = params

  const queryParams: Record<string, string | number> = { page, limit }
  if (search) queryParams.search = search

  const res = await portalApiClient.get<{ data: PublicJobListResponse }>(
    '/api/v1/portal/jobs',
    { params: queryParams }
  )

  return res.data.data
}

/**
 * Fetch a single job with its application form schema.
 * API response shape: { data: PublicJobDetail }
 */
export async function getPublicJobAction(id: string): Promise<PublicJobDetail> {
  const res = await portalApiClient.get<{ data: PublicJobDetail }>(
    `/api/v1/portal/jobs/${id}`
  )

  return res.data.data
}

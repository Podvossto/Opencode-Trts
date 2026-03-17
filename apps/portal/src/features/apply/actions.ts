// Purpose: Portal apply feature — API call functions for application submission
// Ref: API Contract — GET /api/v1/portal/jobs/:id | POST /api/v1/portal/apply
// Role: Public (no auth)

import portalApiClient from '@/lib/api'
import type { PublicJobWithForm } from './types'

/**
 * Fetch a single open job with its form schema.
 * Returns null if job is not found or not open.
 */
export async function getPublicJobAction(jobId: string): Promise<PublicJobWithForm | null> {
  try {
    const res = await portalApiClient.get<{ data: PublicJobWithForm }>(
      `/api/v1/portal/jobs/${jobId}`
    )
    return res.data.data
  } catch (err: unknown) {
    const axiosErr = err as { response?: { status?: number } }
    if (axiosErr?.response?.status === 404) return null
    throw err
  }
}

export interface SubmitApplicationResult {
  id: string
  status: string
}

/**
 * Submit an application with optional resume file.
 * Uses multipart/form-data.
 * Throws with code 'BLACKLISTED' (403) or 'DUPLICATE' (409).
 */
export async function submitApplicationAction(params: {
  jobId: string
  email: string
  firstName: string
  lastName: string
  phone: string
  formData: Record<string, unknown>
  resumeFile?: File
}): Promise<SubmitApplicationResult> {
  const fd = new FormData()
  fd.append('job_id', params.jobId)
  fd.append('email', params.email)
  fd.append('first_name', params.firstName)
  fd.append('last_name', params.lastName)
  fd.append('phone', params.phone)
  fd.append('form_data', JSON.stringify(params.formData))
  if (params.resumeFile) {
    fd.append('resume', params.resumeFile)
  }

  const res = await portalApiClient.post<{ data: SubmitApplicationResult }>(
    '/api/v1/portal/apply',
    fd,
    { headers: { 'Content-Type': 'multipart/form-data' } }
  )
  return res.data.data
}

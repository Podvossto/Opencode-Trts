// Purpose: Axios HTTP client for apps/portal — public requests
// Ref: API Contract — GET /api/v1/portal/jobs
// Portal is public-first: no auth required for job listing

import axios, { type AxiosError } from 'axios'

const portalApiClient = axios.create({
  baseURL: '',  // Empty = relative URLs; Next.js rewrite proxy handles /api/v1/* -> go-api
  headers: { 'Content-Type': 'application/json' },
  timeout: 15000,
})

// Response interceptor — handle errors gracefully
portalApiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    // Let the caller handle the error
    return Promise.reject(error)
  }
)

// Typed helpers
export async function get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  const res = await portalApiClient.get<{ data: T }>(url, { params })
  return res.data.data
}

export async function post<T>(url: string, body?: unknown): Promise<T> {
  const res = await portalApiClient.post<{ data: T }>(url, body)
  return res.data.data
}

export default portalApiClient

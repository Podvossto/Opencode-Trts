// Purpose: Axios HTTP client for apps/main — attaches Bearer JWT, handles 401 redirect
// Ref: API Contract — all authenticated endpoints in /api/v1/*

import axios, { type AxiosError, type AxiosResponse } from 'axios'
import { getToken, clearToken } from '@/lib/auth'

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080',
  headers: { 'Content-Type': 'application/json' },
  timeout: 15000,
})

// Request interceptor — attach Bearer token from cookie
apiClient.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor — handle 401 (expired/invalid JWT)
apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      clearToken()
      // Only redirect if running in browser (not during SSR)
      if (typeof window !== 'undefined') {
        const currentPath = window.location.pathname
        // Don't redirect if already on auth pages
        if (!currentPath.startsWith('/login') && !currentPath.startsWith('/forgot-password')) {
          window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}&expired=1`
        }
      }
    }
    return Promise.reject(error)
  }
)

// Typed helper functions for common HTTP methods
export async function get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  const response = await apiClient.get<T>(url, { params })
  return response.data
}

export async function post<T>(url: string, data?: unknown): Promise<T> {
  const response = await apiClient.post<T>(url, data)
  return response.data
}

export async function put<T>(url: string, data?: unknown): Promise<T> {
  const response = await apiClient.put<T>(url, data)
  return response.data
}

export async function patch<T>(url: string, data?: unknown): Promise<T> {
  const response = await apiClient.patch<T>(url, data)
  return response.data
}

export async function del<T>(url: string): Promise<T> {
  const response = await apiClient.delete<T>(url)
  return response.data
}

export default apiClient

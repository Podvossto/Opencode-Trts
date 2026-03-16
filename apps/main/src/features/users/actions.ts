// Purpose: Users feature — server actions for admin user management CRUD
// Ref: API Contract — GET /api/v1/admin/users | POST /api/v1/admin/users |
//                     GET/PUT/DELETE /api/v1/admin/users/{id} | POST /api/v1/admin/users/{id}/reset-password
//                     GET /api/v1/admin/departments | POST /api/v1/admin/departments
// Role: Admin only

import { get, post, put, del } from '@/lib/api'
import type { User, Department, CreateUserBody, UpdateUserBody } from '@/features/users/types'

// ---------- API Response types ----------

interface ApiSuccess<T> {
  data: T
}

interface PaginatedData {
  data: User[]
  total: number
  page: number
  limit: number
  total_pages: number
}

interface CreateUserResponse {
  data: User
  temporary_password: string
}

// ---------- Users ----------

export async function listUsersAction(
  page = 1,
  limit = 20,
  filters?: { search?: string; role?: string; department_id?: string; is_active?: string }
): Promise<{ data: User[]; total: number; page: number; totalPages: number }> {
  const params = new URLSearchParams()
  params.set('page', String(page))
  params.set('limit', String(limit))
  if (filters?.search) params.set('search', filters.search)
  if (filters?.role) params.set('role', filters.role)
  if (filters?.department_id) params.set('department_id', filters.department_id)
  if (filters?.is_active) params.set('is_active', filters.is_active)

  const res = await get<ApiSuccess<PaginatedData>>(`/api/v1/admin/users?${params.toString()}`)
  return {
    data: res.data.data,
    total: res.data.total,
    page: res.data.page,
    totalPages: res.data.total_pages,
  }
}

export async function getUserAction(id: string): Promise<User> {
  const res = await get<ApiSuccess<User>>(`/api/v1/admin/users/${id}`)
  return res.data
}

export async function createUserAction(
  body: CreateUserBody
): Promise<{ user: User; temporaryPassword: string }> {
  const res = await post<CreateUserResponse>('/api/v1/admin/users', body)
  return { user: res.data, temporaryPassword: res.temporary_password }
}

export async function updateUserAction(
  id: string,
  body: UpdateUserBody
): Promise<User> {
  const res = await put<ApiSuccess<User>>(`/api/v1/admin/users/${id}`, body)
  return res.data
}

export async function deleteUserAction(id: string): Promise<void> {
  await del(`/api/v1/admin/users/${id}`)
}

export async function adminResetPasswordAction(
  id: string,
  newPassword: string
): Promise<void> {
  await post(`/api/v1/admin/users/${id}/reset-password`, {
    new_password: newPassword,
  })
}

// ---------- Departments ----------

export async function listDepartmentsAction(): Promise<Department[]> {
  const res = await get<ApiSuccess<Department[]>>('/api/v1/admin/departments')
  return res.data
}

export async function createDepartmentAction(name: string): Promise<Department> {
  const res = await post<ApiSuccess<Department>>('/api/v1/admin/departments', { name })
  return res.data
}

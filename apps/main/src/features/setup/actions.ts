// Purpose: Server actions for the setup / system configuration feature
// Ref: API Contract — GET/POST/PUT/DELETE /api/v1/setup/*
// Dev must implement: call Go API via lib/api.ts (Admin-only)

'use server';

import type {
  Department,
  UpsertDepartmentRequest,
  EmailTemplate,
  UpdateEmailTemplateRequest,
  SystemSetting,
} from './types';

// --- Departments ---

export async function listDepartments(): Promise<Department[]> {
  // TODO: call api.get('/setup/departments')
  throw new Error('Not implemented');
}

export async function createDepartment(req: UpsertDepartmentRequest): Promise<Department> {
  // TODO: call api.post('/setup/departments', req)
  throw new Error('Not implemented');
}

export async function updateDepartment(
  id: string,
  req: UpsertDepartmentRequest
): Promise<Department> {
  // TODO: call api.put(`/setup/departments/${id}`, req)
  throw new Error('Not implemented');
}

export async function deleteDepartment(id: string): Promise<void> {
  // TODO: call api.delete(`/setup/departments/${id}`)
  // TODO: handle 409 conflict (department in use) and surface to UI
  throw new Error('Not implemented');
}

// --- Email Templates ---

export async function listEmailTemplates(): Promise<EmailTemplate[]> {
  // TODO: call api.get('/setup/email-templates')
  throw new Error('Not implemented');
}

export async function getEmailTemplate(slug: string): Promise<EmailTemplate> {
  // TODO: call api.get(`/setup/email-templates/${slug}`)
  throw new Error('Not implemented');
}

export async function updateEmailTemplate(
  slug: string,
  req: UpdateEmailTemplateRequest
): Promise<EmailTemplate> {
  // TODO: call api.put(`/setup/email-templates/${slug}`, req)
  throw new Error('Not implemented');
}

// --- System Settings ---

export async function getSettings(): Promise<SystemSetting[]> {
  // TODO: call api.get('/setup/settings')
  throw new Error('Not implemented');
}

export async function updateSetting(key: string, value: string): Promise<SystemSetting> {
  // TODO: call api.put(`/setup/settings/${key}`, { value })
  throw new Error('Not implemented');
}

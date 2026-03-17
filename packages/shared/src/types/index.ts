/**
 * Shared TypeScript types for ATS — consumed by apps/main and apps/portal.
 * These types must mirror the Go API response contracts exactly.
 * 
 * Source of truth: ./auth.ts | ./job.ts | ./form.ts (canonical files)
 * This barrel re-exports all canonical types + defines supporting types.
 */

// ---------------------------------------------------------------------------
// Auth — canonical: ./auth.ts
// ---------------------------------------------------------------------------
export type { UserRole, AuthUser, TokenPair, LoginResponse } from './auth';

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: import('./auth').UserRole;
  department_id: string | null;
  is_active: boolean;
  must_change_password: boolean;
  created_at: string; // ISO 8601
}

export interface Department {
  id: string;
  name: string;
  created_at: string;
}

// ---------------------------------------------------------------------------
// Jobs — canonical: ./job.ts
// ---------------------------------------------------------------------------
export type { JobStatus, HardSkillWeight, Job, PublicJob } from './job';

export interface EmploymentType {
  id: string;
  name: string;
}

export interface ExperienceLevel {
  id: string;
  name: string;
}

export interface HardSkill {
  id: string;
  name: string;
}

// ---------------------------------------------------------------------------
// Forms — canonical: ./form.ts
// ---------------------------------------------------------------------------
export type { FormFieldType, FormField, ApplicationForm } from './form';

// ---------------------------------------------------------------------------
// Applications
// ---------------------------------------------------------------------------

export type ApplicationStatus =
  | 'pending'
  | 'in_review'
  | 'testing'
  | 'interviewing'
  | 'hired'
  | 'dropped'
  | 'transferred';

export interface Application {
  id: string;
  candidate_id: string;
  job_id: string;
  status: ApplicationStatus;
  form_data: Record<string, unknown>;
  ocr_status: 'pending' | 'done' | 'failed';
  similarity_score: number | null;
  submitted_at: string;
  created_at: string;
}

export interface Candidate {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone: string;
  is_blacklisted: boolean;
  created_at: string;
}

// ---------------------------------------------------------------------------
// API Standard Response Envelopes
// ---------------------------------------------------------------------------

export interface ApiSuccess<T> {
  data: T;
}

export interface ApiError {
  error: string;
  details?: ValidationErrorDetail[];
}

export interface ValidationErrorDetail {
  field: string;
  message: string;
}

export type ApiResponse<T> = ApiSuccess<T> | ApiError;

// ---------------------------------------------------------------------------
// M06-M15 Re-exports
// ---------------------------------------------------------------------------

export * from './resume-review';
export * from './test';
export * from './interview';
export * from './transfer';
export * from './talent-pool';
export * from './requisition';
export * from './notification';
export * from './setup';

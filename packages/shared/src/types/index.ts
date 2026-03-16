/**
 * Shared TypeScript types for ATS — consumed by apps/main and apps/portal.
 * These types must mirror the Go API response contracts exactly.
 */

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

export type UserRole = 'admin' | 'hr' | 'supervisor';

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  role: UserRole;
  must_change_password: boolean;
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
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
// Jobs
// ---------------------------------------------------------------------------

export type JobStatus = 'draft' | 'open' | 'closed';

export interface Job {
  id: string;
  title: string;
  department: Department;
  employment_type: EmploymentType;
  experience_level: ExperienceLevel;
  form_id: string;
  description: string;
  requirements: string;
  slots: number;
  status: JobStatus;
  publish_at: string | null;
  close_at: string | null;
  hard_skills: JobHardSkill[];
  created_at: string;
}

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

export interface JobHardSkill {
  skill: HardSkill;
  weight: number; // 1-5
}

// ---------------------------------------------------------------------------
// Forms
// ---------------------------------------------------------------------------

export type FieldType =
  | 'text'
  | 'textarea'
  | 'dropdown'
  | 'checkbox'
  | 'radio'
  | 'file_upload'
  | 'date_picker'
  | 'mcq'
  | 'short_answer'
  | 'essay'
  | 'true_false'
  | 'rating_scale'
  | 'image'
  | 'video'
  | 'link';

export interface FormField {
  id: string;
  type: FieldType;
  label: string;
  required: boolean;
  options?: string[];
  order: number;
}

export interface ApplicationForm {
  id: string;
  name: string;
  form_schema: FormField[];
  created_at: string;
}

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

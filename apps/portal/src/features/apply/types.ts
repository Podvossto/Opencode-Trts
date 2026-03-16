// Purpose: Portal apply feature — TypeScript types for application submission
// Ref: API Contract — POST /api/v1/portal/applications, DB tables: candidates, applications
// Dev must implement: SubmitApplicationBody, PublicJobWithForm

export type { FormField, FormFieldType, ApplicationForm } from '@ats/shared';

export interface PublicJobWithForm {
  id: string;
  title: string;
  description: string;
  department_name: string;
  employment_type_name: string;
  experience_level_name: string;
  status: 'open' | 'closed';
  close_at: string | null;
  form_schema: import('@ats/shared').FormField[];
}

export interface SubmitApplicationBody {
  job_id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  form_data: Record<string, unknown>;
  // file uploads handled as multipart before this call
}

export interface AutosaveBody {
  job_id: string;
  candidate_email: string;
  form_data: Record<string, unknown>;
}

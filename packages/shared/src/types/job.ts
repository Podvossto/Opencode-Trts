// Purpose: Shared TypeScript types for Job entity — consumed by apps/main and apps/portal
// Ref: API Contract — GET/POST /api/v1/jobs, /api/v1/public/jobs
// DB table: jobs, job_hard_skills, job_embeddings

export type JobStatus = 'draft' | 'scheduled' | 'open' | 'closed';

export interface HardSkillWeight {
  skill_id: string;
  skill_name?: string;
  weight: 1 | 2 | 3 | 4 | 5;
}

export interface Job {
  id: string;
  title: string;
  description: string;
  status: JobStatus;
  department_id: string;
  employment_type_id: string;
  experience_level_id: string;
  application_form_id?: string;
  hard_skills: HardSkillWeight[];
  publish_at?: string;  // ISO 8601
  close_at?: string;    // ISO 8601
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface PublicJob extends Pick<Job, 'id' | 'title' | 'description' | 'status' | 'publish_at' | 'close_at'> {
  department_name: string;
  employment_type_name: string;
  experience_level_name: string;
}

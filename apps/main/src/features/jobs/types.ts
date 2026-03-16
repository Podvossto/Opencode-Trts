// Purpose: Jobs feature — local types (extends packages/shared)
// Ref: DB tables: jobs, job_hard_skills, departments, employment_types, experience_levels, hard_skills
// Dev must implement: re-export and extend shared Job type

export type { Job, PublicJob, HardSkillWeight, JobStatus } from '@ats/shared';

export interface JobFilters {
  status?: 'draft' | 'open' | 'closed';
  department_id?: string;
  search?: string;
  page?: number;
  limit?: number;
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

export interface CreateJobBody {
  title: string;
  department_id: string;
  employment_type_id: string;
  experience_level_id: string;
  form_id: string;
  description: string;
  requirements?: string;
  slots: number;
  publish_at?: string;
  close_at?: string;
  hard_skills?: Array<{ skill_id: string; weight: 1 | 2 | 3 | 4 | 5 }>;
}

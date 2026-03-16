// Purpose: Shared types for Setup/Admin feature (M14)
// Ref: DB tables: email_templates, jd_templates, departments, employment_types,
//      experience_levels, hard_skills
// Mirrors Go API response contracts

export type EmailTemplateType = 'system' | 'custom';

export interface EmailTemplate {
  id: string;
  code: string;
  name: string;
  type: EmailTemplateType;
  subject: string;
  body_html: string;
  variables: string[];
  is_active: boolean;
  created_at: string;
}

export interface JdTemplate {
  id: string;
  title: string;
  description: string;
  requirements: string | null;
  created_by: string | null;
  created_at: string;
}

export interface MasterDataItem {
  id: string;
  name: string;
  created_at: string;
}

export type MasterDataCategory =
  | 'departments'
  | 'employment-types'
  | 'experience-levels'
  | 'hard-skills';

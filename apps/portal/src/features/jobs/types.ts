// Purpose: Portal jobs feature — TypeScript types for public job listing
// Matches Go backend: PublicJob, PublicJobDetail, PublicJobListResponse

export interface PublicJob {
  id: string
  title: string
  department_name: string
  employment_type_name: string
  experience_level_name: string
  description: string
  requirements?: string | null
  slots: number
  publish_at?: string | null
  close_at?: string | null
}

export interface PublicJobDetail extends PublicJob {
  form_schema: unknown[]
}

export interface PublicJobListResponse {
  data: PublicJob[]
  total: number
  page: number
  limit: number
  total_pages: number
}

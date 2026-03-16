// Purpose: Pipeline feature — TypeScript types for dashboard and candidate management
// Ref: API Contract — GET /api/v1/applications, PUT /api/v1/applications/{id}/status
// Dev must implement: PipelineStats, ApplicationDetail

export type PipelineStatus = 'pending' | 'in_review' | 'testing' | 'interviewing' | 'hired' | 'dropped' | 'transferred';

export interface PipelineStats {
  total: number;
  remaining: number;
  dropped: number;
  hired: number;
  by_status: Record<PipelineStatus, number>;
}

export interface Candidate {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone: string | null;
  is_blacklisted: boolean;
}

export interface ApplicationDocument {
  id: string;
  file_type: string;
  file_url: string;
  file_name: string;
  file_size_kb: number | null;
}

export interface Application {
  id: string;
  candidate: Candidate;
  job_id: string;
  status: PipelineStatus;
  similarity_score: number | null;
  ocr_status: 'pending' | 'done' | 'failed';
  submitted_at: string;
}

export interface ApplicationDetail extends Application {
  form_data: Record<string, unknown>;
  documents: ApplicationDocument[];
}

export interface PipelineFilters {
  status?: PipelineStatus;
  search?: string;
  page?: number;
  limit?: number;
}

// Purpose: TypeScript types for the candidates feature (candidate 360° profile)
// Ref: API Contract — GET /api/v1/candidates/:id, DB table: candidates, applications
// Dev must implement: use these types in candidate profile pages and hooks

export type BlacklistStatus = 'active' | 'blacklisted';

export interface Candidate {
  id: string;
  fullName: string;
  email: string;
  phone: string | null;
  resumeUrl: string | null;
  blacklistStatus: BlacklistStatus;
  blacklistReason: string | null;
  createdAt: string;
}

export interface CandidateApplication {
  id: string;
  jobId: string;
  jobTitle: string;
  status: string;
  pipelineStage: string | null;
  appliedAt: string;
}

export interface CandidateProfile {
  candidate: Candidate;
  applications: CandidateApplication[];
  talentPoolEntry: {
    id: string;
    note: string | null;
    tags: string[];
    addedAt: string;
  } | null;
}

export interface BlacklistRequest {
  reason: string;
}

export interface CandidateListResponse {
  data: Candidate[];
  total: number;
  page: number;
  limit: number;
}

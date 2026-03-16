// Purpose: Shared types for Talent Pool feature (M10)
// Ref: DB tables: candidates (is_starred), candidate_embeddings, job_embeddings
// Mirrors Go API response contracts

export interface StarredCandidate {
  id: string;
  candidate_id: string;
  candidate_name: string;
  candidate_email: string;
  is_starred: boolean;
  is_blacklisted: boolean;
  resume_vector_available: boolean;
  starred_at: string;
}

export interface TalentPoolMatch {
  candidate_id: string;
  candidate_name: string;
  candidate_email: string;
  job_id: string;
  job_title: string;
  similarity_score: number;  // 0.0 - 1.0
  computed_at: string;
}

export interface TalentPoolListResponse {
  data: StarredCandidate[];
  total: number;
  page: number;
  limit: number;
}

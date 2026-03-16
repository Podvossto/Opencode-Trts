// Purpose: TypeScript types for the resume-review feature
// Ref: API Contract — GET/POST /api/v1/resume-review, DB table: resume_reviews
// Dev must implement: use these types in components and API calls

export type ResumeReviewStatus = 'pending' | 'passed' | 'failed' | 'hold';

export interface ResumeReview {
  id: string;
  applicationId: string;
  reviewedById: string;
  status: ResumeReviewStatus;
  score: number | null;
  notes: string | null;
  ocrText: string | null;
  aiScore: number | null;
  createdAt: string;
  updatedAt: string;
}

export interface ApplicationWithResume {
  id: string;
  candidateName: string;
  candidateEmail: string;
  jobTitle: string;
  resumeUrl: string | null;
  resumeReview: ResumeReview | null;
  appliedAt: string;
}

export interface BulkReviewRequest {
  applicationIds: string[];
  status: ResumeReviewStatus;
  notes?: string;
}

export interface ResumeReviewListResponse {
  data: ApplicationWithResume[];
  total: number;
  page: number;
  limit: number;
}

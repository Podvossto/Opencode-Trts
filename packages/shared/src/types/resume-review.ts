// Purpose: Shared types for Resume Review feature (M06)
// Ref: DB tables: position_marks, application_documents (ocr_status, similarity_score)
// Mirrors Go API response contracts

export type OcrStatus = 'pending' | 'processing' | 'done' | 'failed';

export interface PositionMark {
  id: string;
  application_id: string;
  job_id: string;
  job_title: string;
  marked_by: string | null;
  created_at: string;
}

export interface ResumeReviewItem {
  application_id: string;
  candidate_name: string;
  candidate_email: string;
  job_title: string;
  resume_url: string | null;
  ocr_status: OcrStatus;
  ocr_text: string | null;
  similarity_score: number | null;
  position_marks: PositionMark[];
  submitted_at: string;
}

export interface OverallScoreBreakdown {
  jd_similarity: number;      // 0.0 - 1.0
  hard_skill_score: number;   // 0.0 - 1.0
  jd_weight: number;          // 0.0 - 1.0
  skill_weight: number;       // 0.0 - 1.0
  combined_score: number;     // weighted result
}

export type ReviewAction = 'pass' | 'hold' | 'drop' | 'transfer';

export interface BulkReviewPayload {
  application_ids: string[];
  action: ReviewAction;
  notes?: string;
}

// Purpose: Shared types for Requisition feature (M12)
// Ref: DB table: requisitions
// Mirrors Go API response contracts

export type RequisitionStatus = 'draft' | 'submitted' | 'approved' | 'rejected' | 'converted';

export interface Requisition {
  id: string;
  title: string;
  department_id: string | null;
  department_name: string | null;
  employment_type_id: string | null;
  experience_level_id: string | null;
  description: string | null;
  requirements: string | null;
  slots: number;
  status: RequisitionStatus;
  submitted_by: string | null;
  submitted_by_name: string | null;
  reviewed_by: string | null;
  review_note: string | null;
  converted_job_id: string | null;
  submitted_at: string | null;
  reviewed_at: string | null;
  created_at: string;
}

export interface CreateRequisitionPayload {
  title: string;
  department_id?: string;
  employment_type_id?: string;
  experience_level_id?: string;
  description?: string;
  requirements?: string;
  slots: number;
}

export interface ReviewRequisitionPayload {
  status: 'approved' | 'rejected';
  review_note?: string;
}

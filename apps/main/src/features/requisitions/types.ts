// Purpose: TypeScript types for the requisitions feature (headcount approval)
// Ref: API Contract — GET/POST/PUT /api/v1/requisitions, DB table: requisitions
// Dev must implement: use these types in requisitions list and detail components

export type RequisitionStatus =
  | 'draft'
  | 'pending_approval'
  | 'approved'
  | 'rejected'
  | 'fulfilled'
  | 'cancelled';

export interface Requisition {
  id: string;
  jobId: string | null;
  requestedById: string;
  requestedByName: string;
  approvedById: string | null;
  approvedByName: string | null;
  department: string;
  position: string;
  headcount: number;
  justification: string;
  status: RequisitionStatus;
  rejectionNote: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateRequisitionRequest {
  department: string;
  position: string;
  headcount: number;
  justification: string;
}

export interface ReviewRequisitionRequest {
  action: 'approve' | 'reject';
  rejectionNote?: string;
}

export interface RequisitionListResponse {
  data: Requisition[];
  total: number;
  page: number;
  limit: number;
}

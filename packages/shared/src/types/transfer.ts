// Purpose: Shared types for Transfer feature (M09)
// Ref: DB tables: transfer_history
// Mirrors Go API response contracts

export interface TransferHistoryEntry {
  id: string;
  application_id: string;
  from_job_id: string;
  from_job_title: string;
  to_job_id: string;
  to_job_title: string;
  transferred_by: string | null;
  transferred_by_name: string | null;
  reason: string | null;
  created_at: string;
}

export interface InitiateTransferRequest {
  application_id: string;
  to_job_id: string;
  reason?: string;
}

export interface TransferValidation {
  can_transfer: boolean;
  warnings: string[];  // e.g. "Target position is closed"
}

// Purpose: TypeScript types for the transfer feature (inter-department / system transfer)
// Ref: API Contract — POST /api/v1/transfer, DB table: pipeline_actions (action_type: 'transfer')
// Dev must implement: use these types in transfer modal components

export type TransferDestination = 'talent_pool' | 'another_job' | 'external';

export interface TransferRecord {
  id: string;
  applicationId: string;
  candidateName: string;
  fromJobId: string;
  fromJobTitle: string;
  toJobId: string | null;
  toJobTitle: string | null;
  destination: TransferDestination;
  note: string | null;
  transferredById: string;
  createdAt: string;
}

export interface InitiateTransferRequest {
  applicationId: string;
  destination: TransferDestination;
  toJobId?: string;
  note?: string;
}

export interface TransferListResponse {
  data: TransferRecord[];
  total: number;
  page: number;
  limit: number;
}

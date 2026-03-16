// Purpose: Shared types for Notification feature (M13)
// Ref: DB tables: notifications, notification_preferences
// Mirrors Go API response contracts

export type NotificationEventType =
  | 'FRNOT01'  // New application received
  | 'FRNOT02'  // OCR completed
  | 'FRNOT03'  // Test submitted
  | 'FRNOT04'  // Interview score submitted
  | 'FRNOT05'  // Transfer completed
  | 'FRNOT06'  // Requisition submitted
  | 'FRNOT07'  // Requisition approved/rejected
  | 'FRNOT08'  // Job auto-published
  | 'FRNOT09'  // Job auto-closed
  | 'FRNOT10'; // Candidate blacklisted

export interface Notification {
  id: string;
  user_id: string;
  event_type: NotificationEventType;
  payload: Record<string, unknown>;
  is_read: boolean;
  created_at: string;
}

export interface NotificationPreference {
  id: string;
  user_id: string;
  event_type: NotificationEventType;
  in_app: boolean;
  email: boolean;
}

export interface NotificationListResponse {
  data: Notification[];
  unread_count: number;
  total: number;
  page: number;
  limit: number;
}

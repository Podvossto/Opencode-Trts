// Purpose: TypeScript types for the notifications feature
// Ref: API Contract — GET/PUT /api/v1/notifications, DB table: notifications
// Dev must implement: use these types in notification bell and list components

export type NotificationType =
  | 'application_submitted'
  | 'resume_review_done'
  | 'test_completed'
  | 'interview_scheduled'
  | 'interview_result'
  | 'pipeline_action'
  | 'requisition_reviewed'
  | 'system';

export interface Notification {
  id: string;
  userId: string;
  type: NotificationType;
  title: string;
  message: string;
  isRead: boolean;
  referenceId: string | null;
  referenceType: string | null;
  createdAt: string;
}

export interface NotificationListResponse {
  data: Notification[];
  unreadCount: number;
  total: number;
  page: number;
  limit: number;
}

export interface MarkReadRequest {
  notificationIds: string[];
}

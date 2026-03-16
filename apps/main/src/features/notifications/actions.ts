// Purpose: Server actions for the notifications feature
// Ref: API Contract — GET /api/v1/notifications, PUT /api/v1/notifications/read
// Dev must implement: call Go API via lib/api.ts

'use server';

import type { NotificationListResponse, MarkReadRequest } from './types';

/**
 * Fetch the authenticated user's notifications
 * GET /api/v1/notifications?page=&limit=&unread_only=
 * Role: Any authenticated user
 */
export async function listNotifications(
  page = 1,
  limit = 20,
  unreadOnly = false
): Promise<NotificationListResponse> {
  // TODO: call api.get(`/notifications?page=${page}&limit=${limit}&unread_only=${unreadOnly}`)
  throw new Error('Not implemented');
}

/**
 * Mark one or more notifications as read
 * PUT /api/v1/notifications/read
 * Body: MarkReadRequest { notificationIds: string[] }
 * Role: Any authenticated user (own notifications only)
 */
export async function markNotificationsRead(req: MarkReadRequest): Promise<void> {
  // TODO: call api.put('/notifications/read', req)
  throw new Error('Not implemented');
}

/**
 * Mark all unread notifications as read for the current user
 * PUT /api/v1/notifications/read-all
 * Role: Any authenticated user
 */
export async function markAllNotificationsRead(): Promise<void> {
  // TODO: call api.put('/notifications/read-all')
  throw new Error('Not implemented');
}

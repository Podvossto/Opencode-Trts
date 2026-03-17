// Purpose: Dashboard layout shell — sidebar, topnav, auth guard wrapper
// Ref: API Contract — GET /api/v1/auth/me (for user context)
// This file MUST be at (dashboard)/layout.tsx for Next.js to detect it
import type { ReactNode } from 'react'
import { Suspense } from 'react'
import DashboardLayoutClient from './DashboardLayoutClient'

export default function DashboardLayout({ children }: { children: ReactNode }) {
  return (
    <Suspense fallback={<div className="flex min-h-screen items-center justify-center" />}>
      <DashboardLayoutClient>{children}</DashboardLayoutClient>
    </Suspense>
  )
}

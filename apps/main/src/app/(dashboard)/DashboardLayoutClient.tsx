// Purpose: Dashboard layout client component — useSearchParams safely inside Suspense
'use client'

import type { ReactNode } from 'react'
import { useSearchParams } from 'next/navigation'
import { ChangePasswordModal } from '@/features/auth/components/ChangePasswordModal'
import { useSessionTimeout } from '@/hooks/useSessionTimeout'

export default function DashboardLayoutClient({ children }: { children: ReactNode }) {
  const searchParams = useSearchParams()
  const mustChangePassword = searchParams.get('change_password') === '1'

  // ATS-016: Auto-logout when JWT expires (checks every 30s + on tab focus)
  useSessionTimeout()

  return (
    <div className="flex min-h-screen bg-background">
      {/* Force password change modal — blocks interaction until changed */}
      <ChangePasswordModal open={mustChangePassword} />

      {/* Sidebar */}
      <aside className="hidden w-60 shrink-0 border-r border-border bg-white lg:block">
        <div className="flex h-16 items-center border-b border-border px-6">
          <h2 className="text-lg font-medium text-oxford-blue">
            ATS <span className="text-vivid-red">Recruitment</span>
          </h2>
        </div>
        <nav className="space-y-1 p-4">
          {/* TODO: Implement Sidebar nav component with role-based menu */}
          <p className="text-xs text-quick-silver">Navigation placeholder</p>
        </nav>
      </aside>

      {/* Main content area */}
      <div className="flex flex-1 flex-col">
        {/* TopNav */}
        <header className="flex h-16 items-center justify-between border-b border-border bg-white px-6">
          {/* TODO: Implement TopNav with user avatar, notification bell, locale switch */}
          <div />
          <p className="text-xs text-quick-silver">TopNav placeholder</p>
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-auto p-6">
          {children}
        </main>
      </div>
    </div>
  )
}

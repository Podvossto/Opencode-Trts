// Purpose: Auth layout — centered card with TigerSoft branding for login/forgot-password pages
import type { ReactNode } from 'react'

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-oxford-blue-50 via-white to-serene/30 px-4">
      <div className="w-full max-w-md">
        {/* TigerSoft brand header */}
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-medium text-oxford-blue">
            ATS <span className="text-vivid-red">Recruitment</span>
          </h1>
          <p className="mt-2 text-sm text-quick-silver">TigerSoft HR Platform</p>
        </div>
        {children}
      </div>
    </div>
  )
}

import Link from 'next/link'
import { Briefcase } from 'lucide-react'

export default function PublicLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="min-h-screen flex flex-col bg-background">
      {/* Portal Header */}
      <header className="sticky top-0 z-50 w-full border-b border-border bg-white/95 backdrop-blur supports-[backdrop-filter]:bg-white/80">
        <div className="container flex h-16 items-center justify-between px-page-x">
          <Link href="/jobs" className="flex items-center gap-2.5 group">
            <div className="flex h-9 w-9 items-center justify-center rounded-card bg-vivid-red">
              <Briefcase className="h-5 w-5 text-white" />
            </div>
            <div className="flex flex-col">
              <span className="text-base font-medium text-oxford-blue leading-tight">
                TigerSoft
              </span>
              <span className="text-xs text-muted-foreground leading-tight">
                Careers
              </span>
            </div>
          </Link>

          <nav className="flex items-center gap-4">
            <Link
              href="/jobs"
              className="text-sm font-medium text-oxford-blue hover:text-vivid-red transition-colors duration-150"
            >
              Open Positions
            </Link>
          </nav>
        </div>
      </header>

      {/* Main content */}
      <main className="flex-1">
        {children}
      </main>

      {/* Footer */}
      <footer className="border-t border-border bg-oxford-blue text-white">
        <div className="container px-page-x py-8">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <div className="flex h-7 w-7 items-center justify-center rounded-md bg-vivid-red">
                <Briefcase className="h-4 w-4 text-white" />
              </div>
              <span className="text-sm font-medium">TigerSoft Careers</span>
            </div>
            <p className="text-xs text-oxford-blue-200">
              &copy; {new Date().getFullYear()} TigerSoft Co., Ltd. All rights reserved.
            </p>
          </div>
        </div>
      </footer>
    </div>
  )
}

// Purpose: Session timeout hook — polls JWT expiry, auto-logout when token expires
// Ref: ATS-006 — Session timeout via JWT expiry (30min default)
// Usage: Call useSessionTimeout() in DashboardLayout to enable auto-logout
'use client'

import { useEffect, useCallback, useRef } from 'react'
import { useRouter } from 'next/navigation'
import { getToken, isTokenExpired, clearToken } from '@/lib/auth'
import { logoutAction } from '@/features/auth/actions'

// How often to check token expiry (in ms)
const CHECK_INTERVAL_MS = 30_000 // every 30 seconds

// How many ms before expiry to show a warning (optional future enhancement)
// const WARNING_BEFORE_MS = 5 * 60 * 1000 // 5 minutes

export function useSessionTimeout() {
  const router = useRouter()
  const isLoggingOut = useRef(false)

  const performLogout = useCallback(async () => {
    if (isLoggingOut.current) return
    isLoggingOut.current = true

    try {
      await logoutAction()
    } catch {
      // Even if server logout fails, clear client token
      clearToken()
    }

    router.push('/login?reason=session_expired')
  }, [router])

  useEffect(() => {
    // Immediate check on mount
    const token = getToken()
    if (!token || isTokenExpired(token)) {
      performLogout()
      return
    }

    // Periodic check
    const interval = setInterval(() => {
      const currentToken = getToken()
      if (!currentToken || isTokenExpired(currentToken)) {
        performLogout()
      }
    }, CHECK_INTERVAL_MS)

    // Also check on window focus (user returns from another tab)
    function handleVisibilityChange() {
      if (document.visibilityState === 'visible') {
        const currentToken = getToken()
        if (!currentToken || isTokenExpired(currentToken)) {
          performLogout()
        }
      }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange)

    return () => {
      clearInterval(interval)
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    }
  }, [performLogout])
}

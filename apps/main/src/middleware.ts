// Purpose: Next.js route middleware for apps/main — JWT auth guard + role-based redirect
// Next.js requires this file at src/middleware.ts (NOT inside a directory)
// Ref: apps/main/next.config.ts — middleware applies to (dashboard) routes

import type { NextRequest } from 'next/server'
import { NextResponse } from 'next/server'

// Edge-safe JWT payload decode (no verification — edge runtime has no Node crypto)
function decodeJwtPayload(token: string): { sub: string; role: string; exp: number } | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')))
    return payload
  } catch {
    return null
  }
}

export function middleware(request: NextRequest) {
  const token = request.cookies.get('access_token')?.value
  const { pathname } = request.nextUrl

  // No token → redirect to login
  if (!token) {
    const loginUrl = new URL('/login', request.url)
    loginUrl.searchParams.set('redirect', pathname)
    return NextResponse.redirect(loginUrl)
  }

  // Decode and check expiry
  const payload = decodeJwtPayload(token)
  if (!payload || payload.exp * 1000 < Date.now()) {
    // Expired or malformed → clear cookie and redirect
    const loginUrl = new URL('/login', request.url)
    loginUrl.searchParams.set('redirect', pathname)
    const response = NextResponse.redirect(loginUrl)
    response.cookies.delete('access_token')
    return response
  }

  // Role guard: /users/* and /setup/* → admin only
  const adminOnlyPaths = ['/users', '/setup']
  const isAdminRoute = adminOnlyPaths.some(p => pathname.startsWith(p))
  if (isAdminRoute && payload.role !== 'admin') {
    // Non-admin trying to access admin route → redirect to dashboard home
    return NextResponse.redirect(new URL('/jobs', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    // Match all dashboard routes (pages inside (dashboard) group)
    '/jobs/:path*',
    '/pipeline/:path*',
    '/forms/:path*',
    '/resume-review/:path*',
    '/tests/:path*',
    '/interviews/:path*',
    '/candidates/:path*',
    '/requisitions/:path*',
    '/users/:path*',
    '/setup/:path*',
  ],
}

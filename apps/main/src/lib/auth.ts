// Purpose: Auth token helpers for apps/main — read/write/clear JWT from cookies
// Ref: API Contract — POST /api/v1/auth/login (sets access_token), POST /api/v1/auth/logout

import Cookies from 'js-cookie'

export type JWTPayload = {
  sub: string
  role: 'hr' | 'supervisor' | 'admin'
  exp: number
  iat: number
  email?: string
  must_change_password?: boolean
}

const TOKEN_KEY = 'access_token'

/**
 * Read access_token from cookie.
 * Returns the raw JWT string or null if not found.
 */
export function getToken(): string | null {
  return Cookies.get(TOKEN_KEY) ?? null
}

/**
 * Write access_token cookie.
 * @param token - Raw JWT string
 * @param expiresInSeconds - Seconds until expiry (from API response or JWT exp claim)
 */
export function setToken(token: string, expiresInSeconds: number): void {
  // Convert seconds to days for js-cookie
  const days = expiresInSeconds / 86400
  Cookies.set(TOKEN_KEY, token, {
    expires: days,
    path: '/',
    sameSite: 'lax',
    // httpOnly must be false — client needs to read token for role/exp checks
  })
}

/**
 * Remove access_token cookie — used on logout or 401.
 */
export function clearToken(): void {
  Cookies.remove(TOKEN_KEY, { path: '/' })
}

/**
 * Base64-decode JWT payload without cryptographic verification.
 * This is safe for client-side role/exp reads — actual verification
 * happens server-side (Go middleware) and in Next.js edge middleware.
 */
export function decodeToken(token: string): JWTPayload | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null

    // Handle URL-safe base64
    const base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/')
    const jsonStr = atob(base64)
    const payload = JSON.parse(jsonStr) as JWTPayload
    return payload
  } catch {
    return null
  }
}

/**
 * Check if the JWT has expired (based on exp claim).
 * Returns true if expired OR if token is malformed.
 */
export function isTokenExpired(token: string): boolean {
  const payload = decodeToken(token)
  if (!payload) return true
  // exp is in seconds, Date.now() in milliseconds
  return payload.exp * 1000 < Date.now()
}

/**
 * Get the decoded payload from the stored cookie token.
 * Returns null if no token or token is expired.
 */
export function getCurrentUser(): JWTPayload | null {
  const token = getToken()
  if (!token || isTokenExpired(token)) return null
  return decodeToken(token)
}

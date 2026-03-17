// Purpose: Auth feature — actions for login, logout, OTP, password change
// Ref: API Contract — POST /api/v1/auth/login | /logout | /otp/request |
//                     POST /api/v1/auth/password/reset | PUT /api/v1/auth/password/change

import { post, put } from '@/lib/api'
import { setToken, clearToken, getToken, decodeToken } from '@/lib/auth'
import type { LoginResponse } from '@ats/shared'

// ---------- Types ----------

interface ApiSuccess<T> {
  data: T
  message?: string
}

interface LoginResult {
  success: boolean
  mustChangePassword?: boolean
  error?: string
}

interface OtpVerifyResult {
  reset_token: string
}

// ---------- Login ----------

export async function loginAction(email: string, password: string): Promise<LoginResult> {
  try {
    const res = await post<ApiSuccess<LoginResponse>>('/api/v1/auth/login', {
      email,
      password,
    })

    // Go API returns flat fields: access_token, expires_in, role, must_change_password
    const { access_token, expires_in, must_change_password } = res.data

    // Decode token to get expiry for cookie TTL — fall back to API-provided expires_in
    const payload = decodeToken(access_token)
    const expiresIn = payload
      ? payload.exp - Math.floor(Date.now() / 1000)
      : (expires_in ?? 1800)

    setToken(access_token, expiresIn)

    return {
      success: true,
      mustChangePassword: must_change_password,
    }
  } catch (err: unknown) {
    const axiosErr = err as { response?: { status?: number; data?: { message?: string; error?: string } } }
    if (axiosErr.response?.status === 401) {
      return { success: false, error: 'Invalid email or password' }
    }
    if (axiosErr.response?.status === 403) {
      return { success: false, error: 'Account is deactivated. Contact your administrator.' }
    }
    return { success: false, error: 'An unexpected error occurred. Please try again.' }
  }
}

// ---------- Logout ----------

export async function logoutAction(): Promise<void> {
  try {
    const token = getToken()
    if (token) {
      await post('/api/v1/auth/logout')
    }
  } catch {
    // Silently fail — token will be cleared regardless
  } finally {
    clearToken()
  }
}

// ---------- OTP ----------

export async function requestOTPAction(email: string): Promise<{ success: boolean; error?: string }> {
  try {
    await post('/api/v1/auth/otp/request', { email })
    return { success: true }
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { message?: string } } }
    return {
      success: false,
      error: axiosErr.response?.data?.message ?? 'Failed to send OTP',
    }
  }
}

export async function verifyOTPAction(
  email: string,
  otp: string
): Promise<{ success: boolean; resetToken?: string; error?: string }> {
  try {
    const res = await post<ApiSuccess<OtpVerifyResult>>('/api/v1/auth/otp/verify', {
      email,
      otp,
    })
    return { success: true, resetToken: res.data.reset_token }
  } catch (err: unknown) {
    const axiosErr = err as { response?: { status?: number; data?: { message?: string } } }
    if (axiosErr.response?.status === 400) {
      return { success: false, error: 'Invalid or expired OTP' }
    }
    return { success: false, error: 'Verification failed. Please try again.' }
  }
}

// ---------- Password ----------

export async function resetPasswordAction(
  token: string,
  newPassword: string
): Promise<{ success: boolean; error?: string }> {
  try {
    await post('/api/v1/auth/password/reset', {
      token,
      new_password: newPassword,
    })
    return { success: true }
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { message?: string } } }
    return {
      success: false,
      error: axiosErr.response?.data?.message ?? 'Password reset failed',
    }
  }
}

export async function changePasswordAction(
  oldPassword: string,
  newPassword: string
): Promise<{ success: boolean; error?: string }> {
  try {
    await put('/api/v1/auth/password/change', {
      old_password: oldPassword,
      new_password: newPassword,
    })
    return { success: true }
  } catch (err: unknown) {
    const axiosErr = err as { response?: { status?: number; data?: { message?: string } } }
    if (axiosErr.response?.status === 400) {
      return { success: false, error: 'Current password is incorrect' }
    }
    return {
      success: false,
      error: axiosErr.response?.data?.message ?? 'Password change failed',
    }
  }
}

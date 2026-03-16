// Purpose: Shared TypeScript types for auth — consumed by apps/main and services/go-api contract
// Ref: API Contract — POST /api/v1/auth/login, /change-password, /otp/*
// Dev must implement: extend as needed, do NOT use raw strings for role — use UserRole

export type UserRole = 'admin' | 'hr' | 'supervisor';

export interface AuthUser {
  id: string;
  username: string;
  role: UserRole;
  must_change_password: boolean;
  locale: 'th' | 'en';
  theme: 'light' | 'dark';
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export interface LoginResponse extends TokenPair {
  user: AuthUser;
}

// Purpose: Auth feature — local TypeScript types (extends packages/shared/src/types/auth.ts)
// Ref: API Contract — POST /api/v1/auth/login
// Dev must implement: re-export shared types + any feature-specific extensions

export type { UserRole, AuthUser, TokenPair, LoginResponse } from '@ats/shared';

export interface LoginFormValues {
  email: string;
  password: string;
}

export interface ChangePasswordFormValues {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

export interface OTPFormValues {
  email: string;
}

export interface OTPVerifyFormValues {
  email: string;
  otp: string;
}

export interface ResetPasswordFormValues {
  token: string;
  new_password: string;
  confirm_password: string;
}

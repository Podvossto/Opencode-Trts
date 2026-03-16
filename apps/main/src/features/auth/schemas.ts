// Purpose: Auth feature — Zod validation schemas for login, OTP, password forms
// Ref: API Contract — POST /api/v1/auth/login, /otp/request, /otp/verify, /password/reset, /password/change
// Dev must implement: export all schemas; use with react-hook-form + zodResolver

import { z } from 'zod';

export const loginSchema = z.object({
  email: z.string().email('Invalid email'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
});

export const otpRequestSchema = z.object({
  email: z.string().email('Invalid email'),
});

export const otpVerifySchema = z.object({
  email: z.string().email(),
  otp: z.string().length(6, 'OTP must be 6 digits'),
});

export const resetPasswordSchema = z.object({
  token: z.string().min(1),
  new_password: z.string().min(8, 'Password must be at least 8 characters'),
  confirm_password: z.string(),
}).refine(d => d.new_password === d.confirm_password, {
  message: 'Passwords do not match',
  path: ['confirm_password'],
});

export const changePasswordSchema = z.object({
  old_password: z.string().min(8),
  new_password: z.string().min(8, 'Password must be at least 8 characters'),
  confirm_password: z.string(),
}).refine(d => d.new_password === d.confirm_password, {
  message: 'Passwords do not match',
  path: ['confirm_password'],
});

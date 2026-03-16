// Purpose: Users feature — Zod validation schemas for create/edit user forms
// Ref: API Contract — POST /api/v1/admin/users, PUT /api/v1/admin/users/:id

import { z } from 'zod'

export const createUserSchema = z.object({
  email: z.string().email('Invalid email'),
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  role: z.enum(['admin', 'hr', 'supervisor'], { required_error: 'Role is required' }),
  department_id: z.string().optional(),
})

export const updateUserSchema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  role: z.enum(['admin', 'hr', 'supervisor'], { required_error: 'Role is required' }),
  department_id: z.string().optional(),
  is_active: z.boolean().optional(),
})

export const adminResetPasswordSchema = z.object({
  new_password: z.string().min(8, 'Password must be at least 8 characters'),
  confirm_password: z.string(),
}).refine(d => d.new_password === d.confirm_password, {
  message: 'Passwords do not match',
  path: ['confirm_password'],
})

export type CreateUserFormValues = z.infer<typeof createUserSchema>
export type UpdateUserFormValues = z.infer<typeof updateUserSchema>
export type AdminResetPasswordFormValues = z.infer<typeof adminResetPasswordSchema>

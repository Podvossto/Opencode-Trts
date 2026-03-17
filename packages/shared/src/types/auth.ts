// Purpose: Shared TypeScript types for auth — consumed by apps/main and apps/portal
// Ref: API Contract — POST /api/v1/auth/login response shape (Go service)
// Go LoginResponse: { access_token, expires_in, role, must_change_password }

export type UserRole = 'admin' | 'hr' | 'supervisor';

// AuthUser is used for user management (CRUD), NOT for login response
export interface AuthUser {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  is_active: boolean;
  must_change_password: boolean;
  department_id?: string;
  created_at?: string;
}

// TokenPair — base token fields
export interface TokenPair {
  access_token: string;
  expires_in: number; // seconds
}

// LoginResponse — matches Go API flat response exactly
// POST /api/v1/auth/login → { data: LoginResponse }
export interface LoginResponse extends TokenPair {
  role: UserRole;
  must_change_password: boolean;
}

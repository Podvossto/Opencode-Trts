// Purpose: Users feature — TypeScript types for user management
// Ref: API Contract — GET/POST /api/v1/admin/users, DB table: users, departments
// Dev must implement: re-export and extend as needed

export type UserRole = 'admin' | 'hr' | 'supervisor';

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  department_id: string | null;
  department_name?: string;
  is_active: boolean;
  must_change_password: boolean;
  created_at: string;
}

export interface Department {
  id: string;
  name: string;
  created_at: string;
}

export interface CreateUserBody {
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  department_id?: string;
}

export interface UpdateUserBody {
  first_name?: string;
  last_name?: string;
  role?: UserRole;
  department_id?: string | null;
  is_active?: boolean;
}

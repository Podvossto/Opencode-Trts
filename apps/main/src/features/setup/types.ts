// Purpose: TypeScript types for the setup / system configuration feature
// Ref: API Contract — GET/POST/PUT/DELETE /api/v1/setup/*
// Dev must implement: use these types in Admin setup panels

export interface Department {
  id: string;
  name: string;
  code: string;
  managerId: string | null;
  managerName: string | null;
  createdAt: string;
}

export interface UpsertDepartmentRequest {
  name: string;
  code: string;
  managerId?: string;
}

export interface EmailTemplate {
  id: string;
  slug: string;
  subject: string;
  body: string;
  variables: string[];
  updatedAt: string;
}

export interface UpdateEmailTemplateRequest {
  subject: string;
  body: string;
}

export interface SystemSetting {
  key: string;
  value: string;
  updatedAt: string;
}

export interface UpdateSettingRequest {
  value: string;
}

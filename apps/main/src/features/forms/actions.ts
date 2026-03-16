// Purpose: Forms feature — server actions for application form builder CRUD
// Ref: API Contract — GET/POST /api/v1/forms | GET/PUT/DELETE /api/v1/forms/{id}
// DB table: application_forms (form_schema JSONB)
// Role: HR/Admin
// Dev must implement:
//   listFormsAction(): Promise<ApplicationForm[]>
//   getFormAction(id): Promise<ApplicationForm>
//   createFormAction(body): Promise<ApplicationForm>
//   updateFormAction(id, body): Promise<ApplicationForm>
//   deleteFormAction(id): Promise<void>
// IMPORTANT: delete must check for job references — return 409 if form is in use

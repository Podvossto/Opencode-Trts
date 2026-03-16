// Purpose: Portal apply feature — server actions for public job application submission
// Ref: API Contract — GET /api/v1/portal/jobs/{id} | POST /api/v1/portal/applications |
//                     POST /api/v1/portal/applications/autosave
// DB tables: candidates, applications, application_documents, test_autosave
// Role: Public (no auth)
// Dev must implement:
//   getPublicJobAction(jobId): Promise<PublicJobWithForm>  — fetches job + form schema
//   submitApplicationAction(body: SubmitApplicationBody): Promise<{ id: string }>
//   autosaveAction(jobId, email, formData): Promise<void>
//   getAutosaveAction(jobId, email): Promise<Partial<FormData> | null>
// IMPORTANT: submitApplicationAction returns 201 immediately; OCR is async server-side
// BLACKLIST: if candidate email is blacklisted, API returns 403 — show branded error page

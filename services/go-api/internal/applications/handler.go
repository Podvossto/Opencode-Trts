// Purpose: HTTP handlers for the Applications domain — portal submission, autosave, status view
// Ref: API Contract — POST /api/v1/portal/applications, GET /api/v1/portal/applications/{id}
//
//	POST /api/v1/portal/applications/autosave, GET /api/v1/applications
//
// DB tables: candidates, applications, application_documents, test_autosave
// Dev must implement: SubmitApplication, GetApplicationStatus, AutosaveApplication, ListApplications, GetApplication
// IMPORTANT: SubmitApplication returns 201 immediately — OCR is enqueued async via ocr_worker
package applications

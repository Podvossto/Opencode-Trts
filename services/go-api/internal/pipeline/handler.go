// Purpose: HTTP handlers for Pipeline actions — stage transitions, email triggers, hire/drop/transfer
// Ref: API Contract — PUT /api/v1/applications/{id}/status, POST /api/v1/applications/{id}/email
//
//	PUT /api/v1/applications/{id}/transfer
//
// DB tables: applications
// Dev must implement: UpdateApplicationStatus, SendApplicantEmail, TransferApplication
package pipeline

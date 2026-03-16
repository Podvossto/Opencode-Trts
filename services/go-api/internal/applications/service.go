// Purpose: Business logic implementation for the Applications domain
// Ref: API Contract — ApplicationService interface defined in routes.go
// DB tables: candidates, applications, application_documents
// Dev must implement: applicationService struct satisfying ApplicationService interface
// Key: blacklist check (applications.status='dropped'), duplicate check (same candidate+job+open)
package applications

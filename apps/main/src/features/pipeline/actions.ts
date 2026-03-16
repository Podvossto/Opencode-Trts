// Purpose: Pipeline feature — server actions for pipeline dashboard and applicant management
// Ref: API Contract — GET /api/v1/applications | GET /api/v1/applications/{id}
//                     PUT /api/v1/applications/{id}/status
//                     POST /api/v1/applications/{id}/email
//                     PUT /api/v1/applications/{id}/transfer
// DB tables: applications, candidates, application_documents
// Role: HR/Supervisor
// Dev must implement:
//   listApplicationsAction(jobId, filters?): Promise<{ data: Application[]; stats: PipelineStats }>
//   getApplicationAction(id): Promise<ApplicationDetail>
//   updateApplicationStatusAction(id, status, note?): Promise<Application>
//   sendEmailAction(id, subject, body): Promise<void>
//   transferApplicationAction(id, targetJobId): Promise<Application>

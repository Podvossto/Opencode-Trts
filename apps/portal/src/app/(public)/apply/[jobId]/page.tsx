// Purpose: Application submission page — applicant fills dynamic form for a specific job
// Ref: API Contract — GET /api/v1/public/jobs/:id/form, POST /api/v1/applications
// DB table: applications, application_forms, form_fields
// Dev must implement: DynamicFormRenderer from form_fields schema, file upload (resume), submit handler
export default function ApplyPage({ params }: { params: { jobId: string } }) {
  return (
    <div>
      {/* TODO: Implement DynamicFormRenderer for jobId={params.jobId} */}
      <p>Apply Page — Job: {params.jobId}</p>
    </div>
  )
}

// Purpose: Resume review page — AI-scored resumes for a specific job, HR review actions
// Ref: API Contract — GET /api/v1/resume-review/:jobId, POST /api/v1/resume-review/:id/action
// DB table: resume_reviews, applications
// Dev must implement: ResumeReviewList with AI score badges, Pass/Fail/Hold actions, PDF preview drawer
export default function ResumeReviewPage({ params }: { params: { jobId: string } }) {
  return (
    <div>
      {/* TODO: Implement ResumeReviewList for jobId={params.jobId} */}
      <p>Resume Review Page — Job: {params.jobId}</p>
    </div>
  )
}

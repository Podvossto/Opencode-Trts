// Purpose: Candidate detail page — full profile, application history, pipeline actions
// Ref: API Contract — GET /api/v1/candidates/:id, POST /api/v1/pipeline/action
// DB table: applications, pipeline_actions, resume_reviews, interviews, tests
// Dev must implement: CandidateProfile, ApplicationTimeline, ActionPanel (Hire/Drop/Transfer)
export default function CandidateDetailPage({ params }: { params: { candidateId: string } }) {
  return (
    <div>
      {/* TODO: Implement CandidateProfile for candidateId={params.candidateId} */}
      <p>Candidate Detail Page — ID: {params.candidateId}</p>
    </div>
  )
}

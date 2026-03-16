// Purpose: Interview scheduling page — schedule, track, and score interviews per job
// Ref: API Contract — GET /api/v1/interviews/:jobId, POST /api/v1/interviews
// DB table: interviews
// Dev must implement: InterviewList, ScheduleInterviewModal, ScoreForm, status badge (scheduled/done/no-show)
export default function InterviewsPage({ params }: { params: { jobId: string } }) {
  return (
    <div>
      {/* TODO: Implement InterviewList for jobId={params.jobId} */}
      <p>Interviews Page — Job: {params.jobId}</p>
    </div>
  )
}

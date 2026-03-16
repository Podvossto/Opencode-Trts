// Purpose: Online test management page — assign, schedule, view results per job
// Ref: API Contract — GET /api/v1/tests/:jobId, POST /api/v1/tests/:id/assign
// DB table: tests, test_assignments
// Dev must implement: TestList, AssignTestModal, ResultsDrawer with score breakdown
export default function TestsPage({ params }: { params: { jobId: string } }) {
  return (
    <div>
      {/* TODO: Implement TestList for jobId={params.jobId} */}
      <p>Tests Page — Job: {params.jobId}</p>
    </div>
  )
}

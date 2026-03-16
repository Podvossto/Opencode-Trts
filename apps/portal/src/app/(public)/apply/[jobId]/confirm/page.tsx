// Purpose: Application confirmation page — shown after successful submission
// Ref: API Contract — (no further API call — uses data from prior POST /api/v1/applications)
// Dev must implement: ConfirmationCard with application reference number, next steps info
export default function ApplyConfirmPage({ params }: { params: { jobId: string } }) {
  return (
    <div>
      {/* TODO: Implement ConfirmationCard after application submit for jobId={params.jobId} */}
      <p>Application Submitted — Job: {params.jobId}</p>
    </div>
  )
}

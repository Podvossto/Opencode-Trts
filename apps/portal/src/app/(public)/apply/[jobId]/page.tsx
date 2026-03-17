import { notFound } from 'next/navigation'
import { getPublicJobAction } from '@/features/apply/actions'
import ApplicationForm from '@/features/apply/components/ApplicationForm'

interface ApplyPageProps {
  params: { jobId: string }
}

export default async function ApplyPage({ params }: ApplyPageProps) {
  const job = await getPublicJobAction(params.jobId)

  if (!job || job.status !== 'open') {
    notFound()
  }

  return (
    <div className="min-h-screen bg-[#F5F6FA]">
      <div className="max-w-2xl mx-auto px-4 py-10">
        {/* Header */}
        <div className="mb-8">
          <p className="text-xs font-semibold text-[#F4001A] uppercase tracking-wider mb-1">
            {job.department_name}
          </p>
          <h1 className="text-2xl font-bold text-[#0B1F3A] mb-1">{job.title}</h1>
          <div className="flex flex-wrap gap-2 mt-2">
            <span className="text-xs bg-[#0B1F3A]/5 text-[#0B1F3A] rounded-full px-3 py-1">
              {job.employment_type_name}
            </span>
            <span className="text-xs bg-[#0B1F3A]/5 text-[#0B1F3A] rounded-full px-3 py-1">
              {job.experience_level_name}
            </span>
            {job.close_at && (
              <span className="text-xs bg-amber-50 text-amber-700 rounded-full px-3 py-1">
                Closes {new Date(job.close_at).toLocaleDateString()}
              </span>
            )}
          </div>
          {job.description && (
            <p className="mt-4 text-sm text-[#6B7280] leading-relaxed">{job.description}</p>
          )}
        </div>

        {/* Application Form */}
        <ApplicationForm jobId={params.jobId} fields={job.form_schema ?? []} />
      </div>
    </div>
  )
}

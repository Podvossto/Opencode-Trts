import Link from 'next/link'
import { CheckCircle } from 'lucide-react'

interface ApplyConfirmPageProps {
  params: { jobId: string }
  searchParams: { ref?: string }
}

export default function ApplyConfirmPage({ params, searchParams }: ApplyConfirmPageProps) {
  const ref = searchParams.ref

  return (
    <div className="min-h-screen bg-[#F5F6FA] flex items-center justify-center px-4">
      <div className="bg-white rounded-2xl border border-[#E5E7EB] p-10 max-w-md w-full text-center shadow-sm">
        <div className="flex justify-center mb-5">
          <div className="w-16 h-16 rounded-full bg-green-50 flex items-center justify-center">
            <CheckCircle size={32} className="text-green-500" />
          </div>
        </div>

        <h1 className="text-xl font-bold text-[#0B1F3A] mb-2">Application Submitted!</h1>
        <p className="text-sm text-[#6B7280] mb-6 leading-relaxed">
          Thank you for your application. We have received your submission and will review it
          shortly.
        </p>

        {ref && (
          <div className="bg-[#F5F6FA] rounded-lg px-4 py-3 mb-6">
            <p className="text-xs text-[#A3A3A3] mb-0.5">Application Reference</p>
            <p className="text-sm font-mono font-semibold text-[#0B1F3A]">{ref}</p>
          </div>
        )}

        <div className="flex flex-col gap-2 text-sm text-[#6B7280] text-left bg-[#F5F6FA] rounded-lg px-4 py-3 mb-8">
          <p className="font-medium text-[#0B1F3A] mb-1">What happens next?</p>
          <p>1. Our team will review your resume and application.</p>
          <p>2. If shortlisted, you will be contacted via email.</p>
          <p>3. You may be invited to complete an online test or interview.</p>
        </div>

        <Link
          href="/jobs"
          className="inline-block bg-[#0B1F3A] hover:bg-[#0B1F3A]/90 text-white text-sm font-semibold px-6 py-2.5 rounded-lg transition"
        >
          Browse More Jobs
        </Link>
      </div>
    </div>
  )
}

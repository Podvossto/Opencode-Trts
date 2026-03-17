'use client'

import React, { useRef, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Upload, X, AlertCircle, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { submitApplicationAction } from '../actions'
import type { FormField } from '@ats/shared'

interface ApplicationFormProps {
  jobId: string
  fields: FormField[]
}

const ACCEPTED_RESUME_TYPES = [
  'application/pdf',
  'application/msword',
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
]
const MAX_FILE_MB = 5

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function FieldRenderer({
  field,
  value,
  onChange,
  error,
}: {
  field: FormField
  value: unknown
  onChange: (val: unknown) => void
  error?: string
}) {
  const labelCls = 'block text-sm font-medium mb-1.5 text-[#0B1F3A]'
  const inputCls =
    'w-full rounded-lg border border-[#DBE1E1] px-3 py-2 text-sm text-[#0B1F3A] placeholder:text-[#A3A3A3] focus:outline-none focus:ring-2 focus:ring-[#F4001A]/30 focus:border-[#F4001A] transition bg-white'
  const errorCls = 'text-xs text-[#F4001A] mt-1'

  switch (field.type) {
    case 'text':
    case 'short_answer':
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <input
            type="text"
            className={inputCls}
            placeholder={field.placeholder ?? ''}
            value={(value as string) ?? ''}
            onChange={(e) => onChange(e.target.value)}
          />
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )

    case 'textarea':
    case 'essay':
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <textarea
            className={`${inputCls} min-h-[96px] resize-y`}
            placeholder={field.placeholder ?? ''}
            value={(value as string) ?? ''}
            onChange={(e) => onChange(e.target.value)}
          />
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )

    case 'dropdown':
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <select
            className={inputCls}
            value={(value as string) ?? ''}
            onChange={(e) => onChange(e.target.value)}
          >
            <option value="">Select an option</option>
            {(field.options ?? []).map((opt) => (
              <option key={opt} value={opt}>
                {opt}
              </option>
            ))}
          </select>
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )

    case 'radio':
    case 'true_false':
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <div className="flex flex-col gap-2 mt-1">
            {(field.options ?? (field.type === 'true_false' ? ['True', 'False'] : [])).map(
              (opt) => (
                <label key={opt} className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    name={field.field_id}
                    value={opt}
                    checked={(value as string) === opt}
                    onChange={() => onChange(opt)}
                    className="accent-[#F4001A]"
                  />
                  <span className="text-sm text-[#0B1F3A]">{opt}</span>
                </label>
              )
            )}
          </div>
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )

    case 'checkbox':
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <div className="flex flex-col gap-2 mt-1">
            {(field.options ?? []).map((opt) => {
              const checked = Array.isArray(value) && (value as string[]).includes(opt)
              return (
                <label key={opt} className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={checked}
                    onChange={() => {
                      const cur = Array.isArray(value) ? (value as string[]) : []
                      onChange(checked ? cur.filter((v) => v !== opt) : [...cur, opt])
                    }}
                    className="accent-[#F4001A]"
                  />
                  <span className="text-sm text-[#0B1F3A]">{opt}</span>
                </label>
              )
            })}
          </div>
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )

    case 'date_picker':
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <input
            type="date"
            className={inputCls}
            value={(value as string) ?? ''}
            onChange={(e) => onChange(e.target.value)}
          />
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )

    case 'rating_scale': {
      const rating = (value as number) ?? 0
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <div className="flex gap-2 mt-1">
            {[1, 2, 3, 4, 5].map((n) => (
              <button
                key={n}
                type="button"
                onClick={() => onChange(n)}
                className={`w-9 h-9 rounded-lg border text-sm font-medium transition ${
                  rating >= n
                    ? 'bg-[#F4001A] border-[#F4001A] text-white'
                    : 'border-[#DBE1E1] text-[#A3A3A3] hover:border-[#F4001A]'
                }`}
              >
                {n}
              </button>
            ))}
          </div>
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )
    }

    case 'mcq':
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <div className="flex flex-col gap-2 mt-1">
            {(field.options ?? []).map((opt) => (
              <label key={opt} className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name={field.field_id}
                  value={opt}
                  checked={(value as string) === opt}
                  onChange={() => onChange(opt)}
                  className="accent-[#F4001A]"
                />
                <span className="text-sm text-[#0B1F3A]">{opt}</span>
              </label>
            ))}
          </div>
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )

    case 'file_upload':
    case 'image':
    case 'video': {
      const fileVal = value as File | null
      const accept = field.accept ?? (field.type === 'image' ? 'image/*' : field.type === 'video' ? 'video/*' : '*/*')
      const maxMb = field.max_size_mb ?? 10
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <label className="flex items-center gap-3 w-full rounded-lg border border-dashed border-[#DBE1E1] px-4 py-3 cursor-pointer hover:border-[#F4001A] transition bg-white">
            <Upload size={16} className="text-[#A3A3A3] shrink-0" />
            <span className="text-sm text-[#A3A3A3] truncate">
              {fileVal ? fileVal.name : `Upload file (max ${maxMb}MB)`}
            </span>
            <input
              type="file"
              accept={accept}
              className="sr-only"
              onChange={(e) => {
                const f = e.target.files?.[0]
                if (f) onChange(f)
              }}
            />
          </label>
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )
    }

    case 'link_url':
      return (
        <div>
          <label className={labelCls}>
            {field.label}
            {field.required && <span className="text-[#F4001A] ml-0.5">*</span>}
          </label>
          <input
            type="url"
            className={inputCls}
            placeholder={field.placeholder ?? 'https://'}
            value={(value as string) ?? ''}
            onChange={(e) => onChange(e.target.value)}
          />
          {error && <p className={errorCls}>{error}</p>}
        </div>
      )

    default:
      return null
  }
}

export default function ApplicationForm({ jobId, fields }: ApplicationFormProps) {
  const router = useRouter()
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [email, setEmail] = useState('')
  const [phone, setPhone] = useState('')
  const [resumeFile, setResumeFile] = useState<File | null>(null)
  const [resumeError, setResumeError] = useState('')
  const [fieldValues, setFieldValues] = useState<Record<string, unknown>>({})
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [coreErrors, setCoreErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)
  const [blacklisted, setBlacklisted] = useState(false)
  const [duplicate, setDuplicate] = useState(false)
  const [submitError, setSubmitError] = useState('')

  function handleResumeChange(e: React.ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0]
    if (!f) return
    if (!ACCEPTED_RESUME_TYPES.includes(f.type)) {
      setResumeError('Only PDF, DOC, or DOCX files are accepted.')
      setResumeFile(null)
      return
    }
    if (f.size > MAX_FILE_MB * 1024 * 1024) {
      setResumeError(`File must be smaller than ${MAX_FILE_MB}MB.`)
      setResumeFile(null)
      return
    }
    setResumeError('')
    setResumeFile(f)
  }

  function validate(): boolean {
    const ce: Record<string, string> = {}
    if (!firstName.trim()) ce.firstName = 'First name is required.'
    if (!lastName.trim()) ce.lastName = 'Last name is required.'
    if (!email.trim()) ce.email = 'Email is required.'
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) ce.email = 'Enter a valid email address.'
    if (!phone.trim()) ce.phone = 'Phone number is required.'

    const fe: Record<string, string> = {}
    for (const field of fields) {
      if (!field.required) continue
      const val = fieldValues[field.field_id]
      const isEmpty =
        val === undefined ||
        val === null ||
        val === '' ||
        (Array.isArray(val) && (val as unknown[]).length === 0)
      if (isEmpty) fe[field.field_id] = `${field.label} is required.`
    }

    setCoreErrors(ce)
    setFieldErrors(fe)
    return Object.keys(ce).length === 0 && Object.keys(fe).length === 0
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBlacklisted(false)
    setDuplicate(false)
    setSubmitError('')

    if (!validate()) return

    setSubmitting(true)
    try {
      const result = await submitApplicationAction({
        jobId,
        email,
        firstName,
        lastName,
        phone,
        formData: fieldValues,
        resumeFile: resumeFile ?? undefined,
      })
      router.push(`/apply/${jobId}/confirm?ref=${result.id}`)
    } catch (err: unknown) {
      const axiosErr = err as { response?: { status?: number } }
      if (axiosErr?.response?.status === 403) {
        setBlacklisted(true)
      } else if (axiosErr?.response?.status === 409) {
        setDuplicate(true)
      } else {
        setSubmitError('Something went wrong. Please try again.')
      }
    } finally {
      setSubmitting(false)
    }
  }

  const inputCls =
    'w-full rounded-lg border border-[#DBE1E1] px-3 py-2 text-sm text-[#0B1F3A] placeholder:text-[#A3A3A3] focus:outline-none focus:ring-2 focus:ring-[#F4001A]/30 focus:border-[#F4001A] transition bg-white'
  const labelCls = 'block text-sm font-medium mb-1.5 text-[#0B1F3A]'
  const errorCls = 'text-xs text-[#F4001A] mt-1'

  if (blacklisted) {
    return (
      <div className="flex flex-col items-center justify-center py-16 gap-4 text-center">
        <div className="w-12 h-12 rounded-full bg-[#F4001A]/10 flex items-center justify-center">
          <AlertCircle size={24} className="text-[#F4001A]" />
        </div>
        <h2 className="text-lg font-semibold text-[#0B1F3A]">Application Not Accepted</h2>
        <p className="text-sm text-[#6B7280] max-w-sm">
          Your email address has been restricted from applying. Please contact us if you believe
          this is an error.
        </p>
      </div>
    )
  }

  if (duplicate) {
    return (
      <div className="flex flex-col items-center justify-center py-16 gap-4 text-center">
        <div className="w-12 h-12 rounded-full bg-amber-50 flex items-center justify-center">
          <AlertCircle size={24} className="text-amber-500" />
        </div>
        <h2 className="text-lg font-semibold text-[#0B1F3A]">Already Applied</h2>
        <p className="text-sm text-[#6B7280] max-w-sm">
          You have already submitted an application for this position. Only one application per
          position is allowed.
        </p>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} noValidate className="flex flex-col gap-6">
      {/* Personal Information */}
      <section className="bg-white rounded-xl border border-[#E5E7EB] p-6 flex flex-col gap-4">
        <h2 className="text-base font-semibold text-[#0B1F3A]">Personal Information</h2>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className={labelCls}>
              First Name <span className="text-[#F4001A]">*</span>
            </label>
            <input
              type="text"
              className={inputCls}
              placeholder="John"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
            />
            {coreErrors.firstName && <p className={errorCls}>{coreErrors.firstName}</p>}
          </div>
          <div>
            <label className={labelCls}>
              Last Name <span className="text-[#F4001A]">*</span>
            </label>
            <input
              type="text"
              className={inputCls}
              placeholder="Doe"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
            />
            {coreErrors.lastName && <p className={errorCls}>{coreErrors.lastName}</p>}
          </div>
        </div>

        <div>
          <label className={labelCls}>
            Email Address <span className="text-[#F4001A]">*</span>
          </label>
          <input
            type="email"
            className={inputCls}
            placeholder="john.doe@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          {coreErrors.email && <p className={errorCls}>{coreErrors.email}</p>}
        </div>

        <div>
          <label className={labelCls}>
            Phone Number <span className="text-[#F4001A]">*</span>
          </label>
          <input
            type="tel"
            className={inputCls}
            placeholder="+66 8X XXX XXXX"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
          />
          {coreErrors.phone && <p className={errorCls}>{coreErrors.phone}</p>}
        </div>
      </section>

      {/* Resume Upload */}
      <section className="bg-white rounded-xl border border-[#E5E7EB] p-6 flex flex-col gap-3">
        <h2 className="text-base font-semibold text-[#0B1F3A]">Resume / CV</h2>
        <p className="text-xs text-[#6B7280]">PDF, DOC, or DOCX — max {MAX_FILE_MB}MB</p>

        {resumeFile ? (
          <div className="flex items-center gap-3 rounded-lg border border-[#DBE1E1] px-4 py-3 bg-[#F9FAFB]">
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-[#0B1F3A] truncate">{resumeFile.name}</p>
              <p className="text-xs text-[#A3A3A3]">{formatFileSize(resumeFile.size)}</p>
            </div>
            <button
              type="button"
              onClick={() => {
                setResumeFile(null)
                if (fileInputRef.current) fileInputRef.current.value = ''
              }}
              className="p-1 rounded hover:bg-[#F4001A]/10 text-[#A3A3A3] hover:text-[#F4001A] transition"
            >
              <X size={16} />
            </button>
          </div>
        ) : (
          <label className="flex flex-col items-center gap-2 w-full rounded-lg border-2 border-dashed border-[#DBE1E1] px-6 py-8 cursor-pointer hover:border-[#F4001A] transition group">
            <div className="w-10 h-10 rounded-full bg-[#F4001A]/5 group-hover:bg-[#F4001A]/10 flex items-center justify-center transition">
              <Upload size={18} className="text-[#F4001A]" />
            </div>
            <span className="text-sm font-medium text-[#0B1F3A]">Click to upload resume</span>
            <span className="text-xs text-[#A3A3A3]">PDF, DOC, DOCX up to {MAX_FILE_MB}MB</span>
            <input
              ref={fileInputRef}
              type="file"
              accept=".pdf,.doc,.docx,application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
              className="sr-only"
              onChange={handleResumeChange}
            />
          </label>
        )}
        {resumeError && <p className={errorCls}>{resumeError}</p>}
      </section>

      {/* Dynamic Form Fields */}
      {fields.length > 0 && (
        <section className="bg-white rounded-xl border border-[#E5E7EB] p-6 flex flex-col gap-5">
          <h2 className="text-base font-semibold text-[#0B1F3A]">Application Questions</h2>
          {fields.map((field) => (
            <FieldRenderer
              key={field.field_id}
              field={field}
              value={fieldValues[field.field_id]}
              onChange={(val) =>
                setFieldValues((prev) => ({ ...prev, [field.field_id]: val }))
              }
              error={fieldErrors[field.field_id]}
            />
          ))}
        </section>
      )}

      {/* Submit Error */}
      {submitError && (
        <div className="flex items-center gap-2 rounded-lg border border-[#F4001A]/30 bg-[#F4001A]/5 px-4 py-3">
          <AlertCircle size={16} className="text-[#F4001A] shrink-0" />
          <p className="text-sm text-[#F4001A]">{submitError}</p>
        </div>
      )}

      {/* Submit Button */}
      <div className="flex justify-end">
        <Button
          type="submit"
          disabled={submitting}
          className="bg-[#F4001A] hover:bg-[#D10016] text-white font-semibold px-8 py-2.5 rounded-lg transition disabled:opacity-60 flex items-center gap-2"
        >
          {submitting && <Loader2 size={16} className="animate-spin" />}
          {submitting ? 'Submitting…' : 'Submit Application'}
        </Button>
      </div>
    </form>
  )
}

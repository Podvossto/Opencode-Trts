'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import {
  Search,
  MapPin,
  Clock,
  Users,
  Briefcase,
  ChevronLeft,
  ChevronRight,
  AlertCircle,
} from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { Spinner } from '@/components/ui/spinner'
import { listPublicJobsAction } from '@/features/jobs/actions'
import type { PublicJob, PublicJobListResponse } from '@/features/jobs/types'

// ---------- helpers ----------

function formatDate(iso: string | null | undefined): string {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

function isClosingSoon(closeAt: string | null | undefined): boolean {
  if (!closeAt) return false
  const diff = new Date(closeAt).getTime() - Date.now()
  return diff > 0 && diff < 7 * 24 * 60 * 60 * 1000 // within 7 days
}

// ---------- JobCard ----------

function JobCard({ job }: { job: PublicJob }) {
  const closing = isClosingSoon(job.close_at)

  return (
    <Card className="group hover:shadow-lg transition-shadow duration-250 hover:border-vivid-red/30">
      <CardContent className="p-6">
        {/* Header row */}
        <div className="flex items-start justify-between gap-4 mb-3">
          <div className="flex-1 min-w-0">
            <h3 className="text-lg font-medium text-oxford-blue truncate group-hover:text-vivid-red transition-colors duration-150">
              {job.title}
            </h3>
            <p className="text-sm text-muted-foreground mt-0.5">
              {job.department_name}
            </p>
          </div>
          {closing && (
            <Badge variant="destructive" className="shrink-0 text-[10px]">
              Closing Soon
            </Badge>
          )}
        </div>

        {/* Description preview */}
        <p className="text-sm text-oxford-blue-300 line-clamp-2 mb-4 leading-relaxed">
          {job.description}
        </p>

        {/* Meta chips */}
        <div className="flex flex-wrap gap-2 mb-4">
          <div className="inline-flex items-center gap-1 text-xs text-muted-foreground">
            <Briefcase className="h-3.5 w-3.5" />
            <span>{job.employment_type_name}</span>
          </div>
          <div className="inline-flex items-center gap-1 text-xs text-muted-foreground">
            <MapPin className="h-3.5 w-3.5" />
            <span>{job.experience_level_name}</span>
          </div>
          <div className="inline-flex items-center gap-1 text-xs text-muted-foreground">
            <Users className="h-3.5 w-3.5" />
            <span>{job.slots} position{job.slots !== 1 ? 's' : ''}</span>
          </div>
          {job.close_at && (
            <div className="inline-flex items-center gap-1 text-xs text-muted-foreground">
              <Clock className="h-3.5 w-3.5" />
              <span>Closes {formatDate(job.close_at)}</span>
            </div>
          )}
        </div>

        {/* Action row */}
        <div className="flex items-center justify-between pt-3 border-t border-border">
          {job.publish_at && (
            <span className="text-xs text-muted-foreground">
              Posted {formatDate(job.publish_at)}
            </span>
          )}
          <div className="ml-auto flex gap-2">
            <Link href={`/apply/${job.id}`}>
              <Button size="sm">
                Apply Now
              </Button>
            </Link>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

// ---------- Pagination ----------

function Pagination({
  page,
  totalPages,
  onPageChange,
}: {
  page: number
  totalPages: number
  onPageChange: (p: number) => void
}) {
  if (totalPages <= 1) return null

  return (
    <div className="flex items-center justify-center gap-2 mt-8">
      <Button
        variant="outline"
        size="sm"
        disabled={page <= 1}
        onClick={() => onPageChange(page - 1)}
      >
        <ChevronLeft className="h-4 w-4 mr-1" />
        Previous
      </Button>

      <span className="text-sm text-muted-foreground px-3">
        Page {page} of {totalPages}
      </span>

      <Button
        variant="outline"
        size="sm"
        disabled={page >= totalPages}
        onClick={() => onPageChange(page + 1)}
      >
        Next
        <ChevronRight className="h-4 w-4 ml-1" />
      </Button>
    </div>
  )
}

// ---------- Main Page ----------

export default function PublicJobsPage() {
  const [jobs, setJobs] = useState<PublicJob[]>([])
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [total, setTotal] = useState(0)
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const fetchJobs = useCallback(async (p: number, q: string) => {
    setLoading(true)
    setError(null)
    try {
      const result: PublicJobListResponse = await listPublicJobsAction({
        page: p,
        limit: 12,
        search: q || undefined,
      })
      setJobs(result.data ?? [])
      setTotal(result.total)
      setTotalPages(result.total_pages)
      setPage(result.page)
    } catch {
      setError('Failed to load job listings. Please try again.')
    } finally {
      setLoading(false)
    }
  }, [])

  // Initial load
  useEffect(() => {
    fetchJobs(1, '')
  }, [fetchJobs])

  // Debounced search
  const handleSearchChange = (value: string) => {
    setSearch(value)
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => {
      fetchJobs(1, value)
    }, 400)
  }

  const handlePageChange = (p: number) => {
    fetchJobs(p, search)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  return (
    <div className="container px-page-x py-10">
      {/* Hero section */}
      <div className="text-center mb-10">
        <h1 className="text-3xl md:text-4xl font-medium text-oxford-blue mb-2">
          Join Our Team
        </h1>
        <p className="text-muted-foreground text-base max-w-lg mx-auto">
          Find your next opportunity at TigerSoft. We&apos;re building the future of enterprise software.
        </p>
      </div>

      {/* Search bar */}
      <div className="max-w-xl mx-auto mb-8">
        <div className="relative">
          <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search positions by title or keyword..."
            value={search}
            onChange={(e) => handleSearchChange(e.target.value)}
            className="pl-10"
          />
        </div>
      </div>

      {/* Results summary */}
      {!loading && !error && (
        <div className="flex items-center justify-between mb-6">
          <p className="text-sm text-muted-foreground">
            {total} open position{total !== 1 ? 's' : ''}
            {search && (
              <span>
                {' '}matching &ldquo;<span className="font-medium text-oxford-blue">{search}</span>&rdquo;
              </span>
            )}
          </p>
        </div>
      )}

      {/* Loading state */}
      {loading && (
        <div className="flex flex-col items-center justify-center py-20">
          <Spinner size="lg" />
          <p className="text-sm text-muted-foreground mt-3">Loading positions...</p>
        </div>
      )}

      {/* Error state */}
      {error && !loading && (
        <div className="flex flex-col items-center justify-center py-20">
          <div className="flex items-center gap-2 text-vivid-red mb-3">
            <AlertCircle className="h-5 w-5" />
            <span className="text-sm font-medium">{error}</span>
          </div>
          <Button variant="outline" size="sm" onClick={() => fetchJobs(page, search)}>
            Retry
          </Button>
        </div>
      )}

      {/* Empty state */}
      {!loading && !error && jobs.length === 0 && (
        <div className="flex flex-col items-center justify-center py-20">
          <Briefcase className="h-12 w-12 text-muted-foreground/50 mb-4" />
          <h2 className="text-lg font-medium text-oxford-blue mb-1">
            No open positions
          </h2>
          <p className="text-sm text-muted-foreground max-w-sm text-center">
            {search
              ? `No positions match "${search}". Try a different search term.`
              : 'No open positions at the moment. Please check back later.'}
          </p>
          {search && (
            <Button
              variant="outline"
              size="sm"
              className="mt-4"
              onClick={() => {
                setSearch('')
                fetchJobs(1, '')
              }}
            >
              Clear search
            </Button>
          )}
        </div>
      )}

      {/* Job grid */}
      {!loading && !error && jobs.length > 0 && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            {jobs.map((job) => (
              <JobCard key={job.id} job={job} />
            ))}
          </div>

          <Pagination
            page={page}
            totalPages={totalPages}
            onPageChange={handlePageChange}
          />
        </>
      )}
    </div>
  )
}

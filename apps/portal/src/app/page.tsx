import { redirect } from 'next/navigation'

// Portal root → redirect to /jobs
export default function PortalHomePage() {
  redirect('/jobs')
}

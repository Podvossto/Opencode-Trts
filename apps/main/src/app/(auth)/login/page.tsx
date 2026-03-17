// Purpose: Login page — email + password form for HR/Supervisor/Admin
// Ref: API Contract — POST /api/v1/auth/login
import { Suspense } from 'react'
import LoginForm from './LoginForm'
import { Spinner } from '@/components/ui/spinner'

export default function LoginPage() {
  return (
    <Suspense fallback={<div className="flex justify-center"><Spinner /></div>}>
      <LoginForm />
    </Suspense>
  )
}

// Purpose: Forgot password page — 3-step flow: email → OTP verify → new password
// Ref: API Contract — POST /api/v1/auth/otp/request | /otp/verify | /password/reset
'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import Link from 'next/link'
import { ArrowLeft, Eye, EyeOff, Mail, KeyRound, ShieldCheck } from 'lucide-react'

import { otpRequestSchema, otpVerifySchema, resetPasswordSchema } from '@/features/auth/schemas'
import type { OTPFormValues, OTPVerifyFormValues, ResetPasswordFormValues } from '@/features/auth/types'
import { requestOTPAction, verifyOTPAction, resetPasswordAction } from '@/features/auth/actions'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'

type Step = 'email' | 'otp' | 'reset' | 'success'

export default function ForgotPasswordPage() {
  const router = useRouter()
  const [step, setStep] = useState<Step>('email')
  const [email, setEmail] = useState('')
  const [resetToken, setResetToken] = useState('')
  const [serverError, setServerError] = useState<string | null>(null)

  return (
    <Card>
      {step === 'email' && (
        <EmailStep
          serverError={serverError}
          setServerError={setServerError}
          onSuccess={(submittedEmail) => {
            setEmail(submittedEmail)
            setServerError(null)
            setStep('otp')
          }}
        />
      )}
      {step === 'otp' && (
        <OTPStep
          email={email}
          serverError={serverError}
          setServerError={setServerError}
          onSuccess={(token) => {
            setResetToken(token)
            setServerError(null)
            setStep('reset')
          }}
          onBack={() => {
            setServerError(null)
            setStep('email')
          }}
        />
      )}
      {step === 'reset' && (
        <ResetStep
          resetToken={resetToken}
          serverError={serverError}
          setServerError={setServerError}
          onSuccess={() => {
            setServerError(null)
            setStep('success')
          }}
        />
      )}
      {step === 'success' && (
        <SuccessStep onLogin={() => router.push('/login')} />
      )}
    </Card>
  )
}

// ---------- Step 1: Email ----------

function EmailStep({
  serverError,
  setServerError,
  onSuccess,
}: {
  serverError: string | null
  setServerError: (e: string | null) => void
  onSuccess: (email: string) => void
}) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<OTPFormValues>({
    resolver: zodResolver(otpRequestSchema),
    defaultValues: { email: '' },
  })

  async function onSubmit(values: OTPFormValues) {
    setServerError(null)
    const result = await requestOTPAction(values.email)
    if (!result.success) {
      setServerError(result.error ?? 'Failed to send OTP')
      return
    }
    onSuccess(values.email)
  }

  return (
    <>
      <CardHeader className="text-center">
        <div className="mx-auto mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-vivid-red/10">
          <Mail className="h-6 w-6 text-vivid-red" />
        </div>
        <CardTitle>Forgot Password</CardTitle>
        <CardDescription>
          Enter your email address and we&apos;ll send you an OTP to reset your password.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          {serverError && (
            <div className="rounded-input border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
              {serverError}
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="you@company.com"
              autoComplete="email"
              autoFocus
              {...register('email')}
            />
            {errors.email && (
              <p className="text-xs text-destructive">{errors.email.message}</p>
            )}
          </div>

          <Button type="submit" className="w-full" size="lg" disabled={isSubmitting}>
            {isSubmitting ? (
              <>
                <Spinner size="sm" className="mr-2 text-white" />
                Sending OTP...
              </>
            ) : (
              'Send OTP'
            )}
          </Button>

          <div className="text-center">
            <Link
              href="/login"
              className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-oxford-blue transition-colors"
            >
              <ArrowLeft className="h-3 w-3" />
              Back to Sign In
            </Link>
          </div>
        </form>
      </CardContent>
    </>
  )
}

// ---------- Step 2: OTP Verification ----------

function OTPStep({
  email,
  serverError,
  setServerError,
  onSuccess,
  onBack,
}: {
  email: string
  serverError: string | null
  setServerError: (e: string | null) => void
  onSuccess: (resetToken: string) => void
  onBack: () => void
}) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<OTPVerifyFormValues>({
    resolver: zodResolver(otpVerifySchema),
    defaultValues: { email, otp: '' },
  })

  async function onSubmit(values: OTPVerifyFormValues) {
    setServerError(null)
    const result = await verifyOTPAction(values.email, values.otp)
    if (!result.success) {
      setServerError(result.error ?? 'OTP verification failed')
      return
    }
    onSuccess(result.resetToken!)
  }

  async function handleResend() {
    setServerError(null)
    const result = await requestOTPAction(email)
    if (!result.success) {
      setServerError(result.error ?? 'Failed to resend OTP')
    }
  }

  return (
    <>
      <CardHeader className="text-center">
        <div className="mx-auto mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-vivid-red/10">
          <KeyRound className="h-6 w-6 text-vivid-red" />
        </div>
        <CardTitle>Enter OTP</CardTitle>
        <CardDescription>
          We sent a 6-digit code to <span className="font-medium text-oxford-blue">{email}</span>.
          Enter it below to continue.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          {serverError && (
            <div className="rounded-input border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
              {serverError}
            </div>
          )}

          <input type="hidden" {...register('email')} />

          <div className="space-y-2">
            <Label htmlFor="otp">Verification Code</Label>
            <Input
              id="otp"
              type="text"
              inputMode="numeric"
              pattern="[0-9]*"
              maxLength={6}
              placeholder="000000"
              autoComplete="one-time-code"
              autoFocus
              className="text-center text-lg tracking-[0.5em]"
              {...register('otp')}
            />
            {errors.otp && (
              <p className="text-xs text-destructive">{errors.otp.message}</p>
            )}
          </div>

          <Button type="submit" className="w-full" size="lg" disabled={isSubmitting}>
            {isSubmitting ? (
              <>
                <Spinner size="sm" className="mr-2 text-white" />
                Verifying...
              </>
            ) : (
              'Verify OTP'
            )}
          </Button>

          <div className="flex items-center justify-between text-sm">
            <button
              type="button"
              onClick={onBack}
              className="inline-flex items-center gap-1 text-muted-foreground hover:text-oxford-blue transition-colors"
            >
              <ArrowLeft className="h-3 w-3" />
              Change email
            </button>
            <button
              type="button"
              onClick={handleResend}
              className="text-vivid-red hover:underline underline-offset-4"
            >
              Resend code
            </button>
          </div>
        </form>
      </CardContent>
    </>
  )
}

// ---------- Step 3: Reset Password ----------

function ResetStep({
  resetToken,
  serverError,
  setServerError,
  onSuccess,
}: {
  resetToken: string
  serverError: string | null
  setServerError: (e: string | null) => void
  onSuccess: () => void
}) {
  const [showNew, setShowNew] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: { token: resetToken, new_password: '', confirm_password: '' },
  })

  async function onSubmit(values: ResetPasswordFormValues) {
    setServerError(null)
    const result = await resetPasswordAction(values.token, values.new_password)
    if (!result.success) {
      setServerError(result.error ?? 'Password reset failed')
      return
    }
    onSuccess()
  }

  return (
    <>
      <CardHeader className="text-center">
        <div className="mx-auto mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-vivid-red/10">
          <ShieldCheck className="h-6 w-6 text-vivid-red" />
        </div>
        <CardTitle>Set New Password</CardTitle>
        <CardDescription>
          Choose a strong password with at least 8 characters.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          {serverError && (
            <div className="rounded-input border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
              {serverError}
            </div>
          )}

          <input type="hidden" {...register('token')} />

          <div className="space-y-2">
            <Label htmlFor="new_password">New Password</Label>
            <div className="relative">
              <Input
                id="new_password"
                type={showNew ? 'text' : 'password'}
                placeholder="At least 8 characters"
                autoComplete="new-password"
                autoFocus
                className="pr-10"
                {...register('new_password')}
              />
              <button
                type="button"
                className="absolute right-3 top-1/2 -translate-y-1/2 text-quick-silver hover:text-oxford-blue transition-colors"
                onClick={() => setShowNew(prev => !prev)}
                tabIndex={-1}
                aria-label={showNew ? 'Hide password' : 'Show password'}
              >
                {showNew ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
            {errors.new_password && (
              <p className="text-xs text-destructive">{errors.new_password.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="confirm_password">Confirm New Password</Label>
            <Input
              id="confirm_password"
              type="password"
              placeholder="Re-enter new password"
              autoComplete="new-password"
              {...register('confirm_password')}
            />
            {errors.confirm_password && (
              <p className="text-xs text-destructive">{errors.confirm_password.message}</p>
            )}
          </div>

          <Button type="submit" className="w-full" size="lg" disabled={isSubmitting}>
            {isSubmitting ? (
              <>
                <Spinner size="sm" className="mr-2 text-white" />
                Resetting...
              </>
            ) : (
              'Reset Password'
            )}
          </Button>
        </form>
      </CardContent>
    </>
  )
}

// ---------- Step 4: Success ----------

function SuccessStep({ onLogin }: { onLogin: () => void }) {
  return (
    <>
      <CardHeader className="text-center">
        <div className="mx-auto mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-ufo-green/10">
          <ShieldCheck className="h-6 w-6 text-ufo-green" />
        </div>
        <CardTitle>Password Reset Complete</CardTitle>
        <CardDescription>
          Your password has been updated. You can now sign in with your new credentials.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button onClick={onLogin} className="w-full" size="lg">
          Back to Sign In
        </Button>
      </CardContent>
    </>
  )
}

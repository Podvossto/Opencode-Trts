// Purpose: Force password change modal — shown when must_change_password=true after login
// Ref: API Contract — PUT /api/v1/auth/password/change
'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Eye, EyeOff, ShieldCheck } from 'lucide-react'

import { changePasswordSchema } from '@/features/auth/schemas'
import type { ChangePasswordFormValues } from '@/features/auth/types'
import { changePasswordAction, logoutAction } from '@/features/auth/actions'

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'

interface ChangePasswordModalProps {
  open: boolean
}

export function ChangePasswordModal({ open }: ChangePasswordModalProps) {
  const router = useRouter()
  const [serverError, setServerError] = useState<string | null>(null)
  const [showOld, setShowOld] = useState(false)
  const [showNew, setShowNew] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: { old_password: '', new_password: '', confirm_password: '' },
  })

  async function onSubmit(values: ChangePasswordFormValues) {
    setServerError(null)

    const result = await changePasswordAction(values.old_password, values.new_password)

    if (!result.success) {
      setServerError(result.error ?? 'Password change failed')
      return
    }

    // Password changed successfully — redirect to dashboard (remove query param)
    router.replace('/dashboard')
  }

  async function handleLogout() {
    await logoutAction()
    router.push('/login')
  }

  return (
    <Dialog open={open}>
      <DialogContent
        className="sm:max-w-md"
        onInteractOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
      >
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <ShieldCheck className="h-5 w-5 text-vivid-red" />
            Change Your Password
          </DialogTitle>
          <DialogDescription>
            Your administrator requires you to set a new password before continuing.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 pt-2">
          {/* Server error */}
          {serverError && (
            <div className="rounded-input border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
              {serverError}
            </div>
          )}

          {/* Current password */}
          <div className="space-y-2">
            <Label htmlFor="old_password">Current Password</Label>
            <div className="relative">
              <Input
                id="old_password"
                type={showOld ? 'text' : 'password'}
                placeholder="Enter current password"
                autoComplete="current-password"
                className="pr-10"
                {...register('old_password')}
              />
              <button
                type="button"
                className="absolute right-3 top-1/2 -translate-y-1/2 text-quick-silver hover:text-oxford-blue transition-colors"
                onClick={() => setShowOld(prev => !prev)}
                tabIndex={-1}
                aria-label={showOld ? 'Hide password' : 'Show password'}
              >
                {showOld ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
            {errors.old_password && (
              <p className="text-xs text-destructive">{errors.old_password.message}</p>
            )}
          </div>

          {/* New password */}
          <div className="space-y-2">
            <Label htmlFor="new_password">New Password</Label>
            <div className="relative">
              <Input
                id="new_password"
                type={showNew ? 'text' : 'password'}
                placeholder="At least 8 characters"
                autoComplete="new-password"
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

          {/* Confirm password */}
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

          {/* Actions */}
          <div className="flex flex-col gap-2 pt-2">
            <Button type="submit" className="w-full" size="lg" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Spinner size="sm" className="mr-2 text-white" />
                  Updating...
                </>
              ) : (
                'Update Password'
              )}
            </Button>
            <Button
              type="button"
              variant="ghost"
              className="w-full text-muted-foreground"
              onClick={handleLogout}
              disabled={isSubmitting}
            >
              Sign out instead
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}

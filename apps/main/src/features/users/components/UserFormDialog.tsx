// Purpose: Create/Edit user form dialog
// Ref: API Contract — POST /api/v1/admin/users, PUT /api/v1/admin/users/:id
'use client'

import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { createUserSchema, updateUserSchema, type CreateUserFormValues, type UpdateUserFormValues } from '@/features/users/schemas'
import type { User, Department } from '@/features/users/types'

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

interface UserFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  user: User | null // null = create mode, user = edit mode
  departments: Department[]
  isSubmitting: boolean
  onSubmit: (values: CreateUserFormValues | UpdateUserFormValues) => void
}

export function UserFormDialog({
  open,
  onOpenChange,
  user,
  departments,
  isSubmitting,
  onSubmit,
}: UserFormDialogProps) {
  const isEdit = !!user

  const createForm = useForm<CreateUserFormValues>({
    resolver: zodResolver(createUserSchema),
    defaultValues: { email: '', first_name: '', last_name: '', role: 'hr', department_id: '' },
  })

  const editForm = useForm<UpdateUserFormValues>({
    resolver: zodResolver(updateUserSchema),
    defaultValues: { first_name: '', last_name: '', role: 'hr', department_id: '', is_active: true },
  })

  // Reset form when user changes
  useEffect(() => {
    if (user) {
      editForm.reset({
        first_name: user.first_name,
        last_name: user.last_name,
        role: user.role as 'admin' | 'hr' | 'supervisor',
        department_id: user.department_id ?? '',
        is_active: user.is_active,
      })
    } else {
      createForm.reset({ email: '', first_name: '', last_name: '', role: 'hr', department_id: '' })
    }
  }, [user, createForm, editForm])

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (isEdit) {
      editForm.handleSubmit(onSubmit)(e)
    } else {
      createForm.handleSubmit(onSubmit)(e)
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const form = (isEdit ? editForm : createForm) as ReturnType<typeof useForm<any>>
  const errors = form.formState.errors

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Edit User' : 'Create New User'}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? `Update details for ${user!.first_name} ${user!.last_name}`
              : 'New user will receive a temporary password and must change it on first login.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Email — only shown on create */}
          {!isEdit && (
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="user@company.com"
                {...createForm.register('email')}
              />
              {(errors as Record<string, { message?: string }>).email && (
                <p className="text-xs text-destructive">{(errors as Record<string, { message?: string }>).email?.message}</p>
              )}
            </div>
          )}

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="first_name">First Name</Label>
              <Input
                id="first_name"
                placeholder="First name"
                {...form.register('first_name')}
              />
              {errors.first_name && (
                <p className="text-xs text-destructive">{String(errors.first_name?.message ?? '')}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="last_name">Last Name</Label>
              <Input
                id="last_name"
                placeholder="Last name"
                {...form.register('last_name')}
              />
              {errors.last_name && (
                <p className="text-xs text-destructive">{String(errors.last_name?.message ?? '')}</p>
              )}
            </div>
          </div>

          {/* Role */}
          <div className="space-y-2">
            <Label>Role</Label>
            <Select
              value={form.watch('role')}
              onValueChange={(v) => form.setValue('role', v as 'admin' | 'hr' | 'supervisor')}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select role" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="hr">HR</SelectItem>
                <SelectItem value="supervisor">Supervisor</SelectItem>
                <SelectItem value="admin">Admin</SelectItem>
              </SelectContent>
            </Select>
            {errors.role && (
              <p className="text-xs text-destructive">{String(errors.role?.message ?? '')}</p>
            )}
          </div>

          {/* Department */}
          <div className="space-y-2">
            <Label>Department</Label>
            <Select
              value={form.watch('department_id') ?? ''}
              onValueChange={(v) => form.setValue('department_id', v === '__none__' ? '' : v)}
            >
              <SelectTrigger>
                <SelectValue placeholder="No department" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="__none__">No department</SelectItem>
                {departments.map((d) => (
                  <SelectItem key={d.id} value={d.id}>{d.name}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <DialogFooter className="pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Spinner size="sm" className="mr-2 text-white" />
                  {isEdit ? 'Saving...' : 'Creating...'}
                </>
              ) : (
                isEdit ? 'Save Changes' : 'Create User'
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

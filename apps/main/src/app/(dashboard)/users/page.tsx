// Purpose: User management page — list/create/edit/deactivate users (admin only)
// Ref: API Contract — GET/POST /api/v1/admin/users, PUT/DELETE /api/v1/admin/users/:id
'use client'

import { useCallback, useEffect, useState } from 'react'
import {
  Search,
  Plus,
  MoreHorizontal,
  Pencil,
  KeyRound,
  UserX,
  ChevronLeft,
  ChevronRight,
  Copy,
  Check,
} from 'lucide-react'

import type { User, Department, CreateUserBody, UpdateUserBody } from '@/features/users/types'
import type { CreateUserFormValues, UpdateUserFormValues } from '@/features/users/schemas'
import {
  listUsersAction,
  createUserAction,
  updateUserAction,
  deleteUserAction,
  adminResetPasswordAction,
  listDepartmentsAction,
} from '@/features/users/actions'
import { UserFormDialog } from '@/features/users/components/UserFormDialog'
import { ResetPasswordDialog } from '@/features/users/components/ResetPasswordDialog'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Spinner } from '@/components/ui/spinner'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

const ROLES = [
  { value: '', label: 'All Roles' },
  { value: 'admin', label: 'Admin' },
  { value: 'hr', label: 'HR' },
  { value: 'supervisor', label: 'Supervisor' },
]

const STATUS_OPTIONS = [
  { value: '', label: 'All Status' },
  { value: 'true', label: 'Active' },
  { value: 'false', label: 'Inactive' },
]

const ROLE_BADGE: Record<string, string> = {
  admin: 'bg-vivid-red/10 text-vivid-red border-vivid-red/20',
  hr: 'bg-blue-50 text-blue-700 border-blue-200',
  supervisor: 'bg-amber-50 text-amber-700 border-amber-200',
}

export default function UsersPage() {
  // Data state
  const [users, setUsers] = useState<User[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [total, setTotal] = useState(0)
  const [totalPages, setTotalPages] = useState(0)
  const [loading, setLoading] = useState(true)

  // Filter state
  const [search, setSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const limit = 20

  // Dialog state
  const [formOpen, setFormOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [formSubmitting, setFormSubmitting] = useState(false)

  const [resetOpen, setResetOpen] = useState(false)
  const [resetUser, setResetUser] = useState<User | null>(null)
  const [resetSubmitting, setResetSubmitting] = useState(false)

  const [deactivateOpen, setDeactivateOpen] = useState(false)
  const [deactivateUser, setDeactivateUser] = useState<User | null>(null)

  // Temp password display after create
  const [tempPasswordInfo, setTempPasswordInfo] = useState<{ email: string; password: string } | null>(null)
  const [copied, setCopied] = useState(false)

  const [error, setError] = useState<string | null>(null)

  // Fetch users
  const fetchUsers = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const result = await listUsersAction(page, limit, {
        search: search || undefined,
        role: roleFilter || undefined,
        is_active: statusFilter || undefined,
      })
      setUsers(result.data)
      setTotal(result.total)
      setTotalPages(result.totalPages)
    } catch {
      setError('Failed to load users')
    } finally {
      setLoading(false)
    }
  }, [page, search, roleFilter, statusFilter])

  // Fetch departments once
  useEffect(() => {
    listDepartmentsAction().then(setDepartments).catch(() => {})
  }, [])

  // Fetch users on filter change
  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  // Debounced search
  const [searchInput, setSearchInput] = useState('')
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearch(searchInput)
      setPage(1)
    }, 400)
    return () => clearTimeout(timer)
  }, [searchInput])

  // ---------- Handlers ----------

  async function handleCreateOrUpdate(values: CreateUserFormValues | UpdateUserFormValues) {
    setFormSubmitting(true)
    try {
      if (editingUser) {
        // Edit mode
        const body: UpdateUserBody = {
          first_name: values.first_name,
          last_name: values.last_name,
          role: values.role as 'admin' | 'hr' | 'supervisor',
          department_id: values.department_id || null,
        }
        if ('is_active' in values) body.is_active = values.is_active
        await updateUserAction(editingUser.id, body)
      } else {
        // Create mode
        const createValues = values as CreateUserFormValues
        const body: CreateUserBody = {
          email: createValues.email,
          first_name: createValues.first_name,
          last_name: createValues.last_name,
          role: createValues.role as 'admin' | 'hr' | 'supervisor',
          department_id: createValues.department_id || undefined,
        }
        const result = await createUserAction(body)
        setTempPasswordInfo({ email: result.user.email, password: result.temporaryPassword })
      }
      setFormOpen(false)
      setEditingUser(null)
      fetchUsers()
    } catch {
      // Error handled within the dialog or as toast in the future
    } finally {
      setFormSubmitting(false)
    }
  }

  async function handleResetPassword(password: string) {
    if (!resetUser) return
    setResetSubmitting(true)
    try {
      await adminResetPasswordAction(resetUser.id, password)
      setResetOpen(false)
      setResetUser(null)
    } catch {
      // Could add toast notification here
    } finally {
      setResetSubmitting(false)
    }
  }

  async function handleDeactivate() {
    if (!deactivateUser) return
    try {
      await deleteUserAction(deactivateUser.id)
      setDeactivateOpen(false)
      setDeactivateUser(null)
      fetchUsers()
    } catch {
      // Could add toast notification here
    }
  }

  function handleCopyPassword() {
    if (tempPasswordInfo) {
      navigator.clipboard.writeText(tempPasswordInfo.password)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-medium text-oxford-blue">User Management</h1>
          <p className="text-sm text-muted-foreground">
            Manage HR, Supervisor, and Admin accounts
          </p>
        </div>
        <Button onClick={() => { setEditingUser(null); setFormOpen(true) }}>
          <Plus className="mr-2 h-4 w-4" />
          Create User
        </Button>
      </div>

      {/* Temp password info banner */}
      {tempPasswordInfo && (
        <Card className="border-ufo-green/30 bg-ufo-green/5">
          <CardContent className="flex items-center justify-between py-4">
            <div>
              <p className="text-sm font-medium text-oxford-blue">
                User created: {tempPasswordInfo.email}
              </p>
              <p className="mt-1 text-sm text-muted-foreground">
                Temporary password: <code className="rounded bg-muted px-2 py-0.5 font-mono text-sm">{tempPasswordInfo.password}</code>
              </p>
            </div>
            <div className="flex gap-2">
              <Button size="sm" variant="outline" onClick={handleCopyPassword}>
                {copied ? <Check className="mr-1 h-3 w-3" /> : <Copy className="mr-1 h-3 w-3" />}
                {copied ? 'Copied' : 'Copy'}
              </Button>
              <Button size="sm" variant="ghost" onClick={() => setTempPasswordInfo(null)}>
                Dismiss
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Filters */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Filters</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap items-center gap-3">
            <div className="relative w-72">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search by name or email..."
                className="pl-9"
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
              />
            </div>
            <Select value={roleFilter} onValueChange={(v) => { setRoleFilter(v === '__all__' ? '' : v); setPage(1) }}>
              <SelectTrigger className="w-40">
                <SelectValue placeholder="All Roles" />
              </SelectTrigger>
              <SelectContent>
                {ROLES.map((r) => (
                  <SelectItem key={r.value || '__all__'} value={r.value || '__all__'}>
                    {r.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v === '__all__' ? '' : v); setPage(1) }}>
              <SelectTrigger className="w-40">
                <SelectValue placeholder="All Status" />
              </SelectTrigger>
              <SelectContent>
                {STATUS_OPTIONS.map((s) => (
                  <SelectItem key={s.value || '__all__'} value={s.value || '__all__'}>
                    {s.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <p className="ml-auto text-sm text-muted-foreground">
              {total} user{total !== 1 ? 's' : ''}
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="flex items-center justify-center py-16">
              <Spinner size="lg" />
            </div>
          ) : error ? (
            <div className="py-16 text-center text-sm text-destructive">{error}</div>
          ) : users.length === 0 ? (
            <div className="py-16 text-center text-sm text-muted-foreground">
              No users found matching your filters.
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Role</TableHead>
                  <TableHead>Department</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-12" />
                </TableRow>
              </TableHeader>
              <TableBody>
                {users.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell className="font-medium text-oxford-blue">
                      {user.first_name} {user.last_name}
                    </TableCell>
                    <TableCell className="text-muted-foreground">{user.email}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className={ROLE_BADGE[user.role] ?? ''}>
                        {user.role.toUpperCase()}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {user.department_name ?? '—'}
                    </TableCell>
                    <TableCell>
                      {user.is_active ? (
                        <Badge variant="outline" className="bg-ufo-green/10 text-ufo-green border-ufo-green/20">
                          Active
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="bg-muted text-muted-foreground">
                          Inactive
                        </Badge>
                      )}
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => { setEditingUser(user); setFormOpen(true) }}>
                            <Pencil className="mr-2 h-4 w-4" />
                            Edit
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => { setResetUser(user); setResetOpen(true) }}>
                            <KeyRound className="mr-2 h-4 w-4" />
                            Reset Password
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          {user.is_active && (
                            <DropdownMenuItem
                              className="text-destructive focus:text-destructive"
                              onClick={() => { setDeactivateUser(user); setDeactivateOpen(true) }}
                            >
                              <UserX className="mr-2 h-4 w-4" />
                              Deactivate
                            </DropdownMenuItem>
                          )}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between border-t border-border px-6 py-4">
            <p className="text-sm text-muted-foreground">
              Page {page} of {totalPages}
            </p>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(Math.max(1, page - 1))}
                disabled={page <= 1}
              >
                <ChevronLeft className="mr-1 h-4 w-4" />
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(Math.min(totalPages, page + 1))}
                disabled={page >= totalPages}
              >
                Next
                <ChevronRight className="ml-1 h-4 w-4" />
              </Button>
            </div>
          </div>
        )}
      </Card>

      {/* Create/Edit Dialog */}
      <UserFormDialog
        open={formOpen}
        onOpenChange={(open) => { setFormOpen(open); if (!open) setEditingUser(null) }}
        user={editingUser}
        departments={departments}
        isSubmitting={formSubmitting}
        onSubmit={handleCreateOrUpdate}
      />

      {/* Reset Password Dialog */}
      <ResetPasswordDialog
        open={resetOpen}
        onOpenChange={(open) => { setResetOpen(open); if (!open) setResetUser(null) }}
        userName={resetUser ? `${resetUser.first_name} ${resetUser.last_name}` : ''}
        isSubmitting={resetSubmitting}
        onSubmit={handleResetPassword}
      />

      {/* Deactivate Confirmation */}
      <AlertDialog open={deactivateOpen} onOpenChange={setDeactivateOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Deactivate User</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to deactivate{' '}
              <span className="font-medium text-oxford-blue">
                {deactivateUser?.first_name} {deactivateUser?.last_name}
              </span>
              ? They will no longer be able to log in. This action can be reversed by an admin.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              onClick={handleDeactivate}
            >
              Deactivate
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}

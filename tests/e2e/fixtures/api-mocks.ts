// Shared API mock helpers for Playwright E2E tests
// Intercepts API calls and returns controlled responses so tests are deterministic.

import { type Page, type Route } from '@playwright/test'

// ---------- Constants ----------

const API_BASE = 'http://localhost:8080'

// Fake JWT — structurally valid (3 base64 parts) with embedded claims.
// NOT cryptographically signed — only used for cookie-based auth bypass in tests.
function makeFakeJwt(claims: {
  sub: string
  role: string
  exp?: number
  jti?: string
  must_change_password?: boolean
}): string {
  const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
  const payload = btoa(
    JSON.stringify({
      sub: claims.sub,
      role: claims.role,
      jti: claims.jti ?? 'test-jti-001',
      exp: claims.exp ?? Math.floor(Date.now() / 1000) + 3600, // 1hr from now
      must_change_password: claims.must_change_password ?? false,
    })
  )
  const signature = btoa('test-signature')
  return `${header}.${payload}.${signature}`
}

// ---------- Test user fixtures ----------

export const TEST_ADMIN = {
  id: '00000000-0000-0000-0000-000000000001',
  email: 'admin@tigersoft.com',
  first_name: 'Admin',
  last_name: 'User',
  role: 'admin' as const,
  is_active: true,
  department_id: null,
  department_name: null,
  must_change_password: false,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
}

export const TEST_HR = {
  id: '00000000-0000-0000-0000-000000000002',
  email: 'hr@tigersoft.com',
  first_name: 'HR',
  last_name: 'Manager',
  role: 'hr' as const,
  is_active: true,
  department_id: '00000000-0000-0000-0000-000000000010',
  department_name: 'Human Resources',
  must_change_password: false,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
}

export const TEST_FORCED_CHANGE = {
  ...TEST_ADMIN,
  id: '00000000-0000-0000-0000-000000000003',
  email: 'newuser@tigersoft.com',
  first_name: 'New',
  last_name: 'Employee',
  must_change_password: true,
}

export const ADMIN_JWT = makeFakeJwt({ sub: TEST_ADMIN.id, role: 'admin' })
export const HR_JWT = makeFakeJwt({ sub: TEST_HR.id, role: 'hr' })
export const FORCED_CHANGE_JWT = makeFakeJwt({
  sub: TEST_FORCED_CHANGE.id,
  role: 'admin',
  must_change_password: true,
})
export const EXPIRED_JWT = makeFakeJwt({
  sub: TEST_ADMIN.id,
  role: 'admin',
  exp: Math.floor(Date.now() / 1000) - 600, // expired 10min ago
})

// ---------- Departments ----------

export const TEST_DEPARTMENTS = [
  { id: '00000000-0000-0000-0000-000000000010', name: 'Human Resources', created_at: '2026-01-01T00:00:00Z' },
  { id: '00000000-0000-0000-0000-000000000011', name: 'Engineering', created_at: '2026-01-01T00:00:00Z' },
  { id: '00000000-0000-0000-0000-000000000012', name: 'Finance', created_at: '2026-01-01T00:00:00Z' },
]

// ---------- Users list ----------

export const TEST_USERS_LIST = [
  TEST_ADMIN,
  TEST_HR,
  {
    id: '00000000-0000-0000-0000-000000000004',
    email: 'supervisor@tigersoft.com',
    first_name: 'Super',
    last_name: 'Visor',
    role: 'supervisor',
    is_active: true,
    department_id: '00000000-0000-0000-0000-000000000011',
    department_name: 'Engineering',
    must_change_password: false,
    created_at: '2026-02-01T00:00:00Z',
    updated_at: '2026-02-01T00:00:00Z',
  },
  {
    id: '00000000-0000-0000-0000-000000000005',
    email: 'disabled@tigersoft.com',
    first_name: 'Disabled',
    last_name: 'Account',
    role: 'hr',
    is_active: false,
    department_id: null,
    department_name: null,
    must_change_password: false,
    created_at: '2026-02-15T00:00:00Z',
    updated_at: '2026-03-01T00:00:00Z',
  },
]

// ---------- Jobs (portal) ----------

export const TEST_JOBS = [
  {
    id: '00000000-0000-0000-0000-000000000101',
    title: 'Senior Go Developer',
    description: 'We are looking for an experienced Go developer to join our backend team. You will work on high-performance APIs serving enterprise clients.',
    department_name: 'Engineering',
    employment_type_name: 'Full-time',
    experience_level_name: 'Senior',
    slots: 2,
    publish_at: '2026-03-01T00:00:00Z',
    close_at: '2026-04-30T00:00:00Z',
  },
  {
    id: '00000000-0000-0000-0000-000000000102',
    title: 'Frontend Engineer (React)',
    description: 'Join our frontend team to build modern, accessible UIs using React and Next.js. Experience with TypeScript and Tailwind CSS is a plus.',
    department_name: 'Engineering',
    employment_type_name: 'Full-time',
    experience_level_name: 'Mid-level',
    slots: 3,
    publish_at: '2026-03-05T00:00:00Z',
    close_at: '2026-05-15T00:00:00Z',
  },
  {
    id: '00000000-0000-0000-0000-000000000103',
    title: 'HR Coordinator',
    description: 'Support the HR team with recruitment coordination, onboarding processes, and employee engagement programs.',
    department_name: 'Human Resources',
    employment_type_name: 'Full-time',
    experience_level_name: 'Junior',
    slots: 1,
    publish_at: '2026-03-10T00:00:00Z',
    close_at: '2026-03-20T00:00:00Z', // closing soon
  },
  {
    id: '00000000-0000-0000-0000-000000000104',
    title: 'Data Analyst Intern',
    description: 'A 6-month internship opportunity to work with real business data and learn analytics tools.',
    department_name: 'Finance',
    employment_type_name: 'Internship',
    experience_level_name: 'Entry-level',
    slots: 2,
    publish_at: '2026-03-12T00:00:00Z',
    close_at: null,
  },
]

// ---------- API Mock Helpers ----------

/**
 * Set up auth cookie to bypass login flow in tests.
 * Use this before navigating to authenticated routes.
 */
export async function loginAsAdmin(page: Page): Promise<void> {
  await page.context().addCookies([
    {
      name: 'access_token',
      value: ADMIN_JWT,
      domain: 'localhost',
      path: '/',
      httpOnly: false,
      secure: false,
      sameSite: 'Lax',
    },
  ])
}

export async function loginAsHR(page: Page): Promise<void> {
  await page.context().addCookies([
    {
      name: 'access_token',
      value: HR_JWT,
      domain: 'localhost',
      path: '/',
      httpOnly: false,
      secure: false,
      sameSite: 'Lax',
    },
  ])
}

/**
 * Mock the login API endpoint.
 * mode: 'success' | 'invalid' | 'disabled' | 'must_change'
 */
export async function mockLoginApi(
  page: Page,
  mode: 'success' | 'invalid' | 'disabled' | 'must_change' = 'success'
): Promise<void> {
  await page.route(`${API_BASE}/api/v1/auth/login`, async (route: Route) => {
    if (route.request().method() !== 'POST') {
      await route.fallback()
      return
    }

    switch (mode) {
      case 'success':
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: {
              access_token: ADMIN_JWT,
              user: TEST_ADMIN,
            },
          }),
        })
        break
      case 'invalid':
        await route.fulfill({
          status: 401,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'invalid email or password' }),
        })
        break
      case 'disabled':
        await route.fulfill({
          status: 403,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'account is deactivated' }),
        })
        break
      case 'must_change':
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: {
              access_token: FORCED_CHANGE_JWT,
              user: TEST_FORCED_CHANGE,
            },
          }),
        })
        break
    }
  })
}

/**
 * Mock the logout API endpoint.
 */
export async function mockLogoutApi(page: Page): Promise<void> {
  await page.route(`${API_BASE}/api/v1/auth/logout`, async (route: Route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ message: 'logged out' }),
    })
  })
}

/**
 * Mock OTP request/verify/reset endpoints for forgot password flow.
 */
export async function mockOtpFlow(page: Page): Promise<void> {
  // OTP request
  await page.route(`${API_BASE}/api/v1/auth/otp/request`, async (route: Route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ message: 'OTP sent' }),
    })
  })

  // OTP verify
  await page.route(`${API_BASE}/api/v1/auth/otp/verify`, async (route: Route) => {
    const body = route.request().postDataJSON()
    if (body?.otp === '123456') {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: { reset_token: 'mock-reset-token-abc' } }),
      })
    } else {
      await route.fulfill({
        status: 400,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'invalid or expired OTP' }),
      })
    }
  })

  // Password reset
  await page.route(`${API_BASE}/api/v1/auth/password/reset`, async (route: Route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ message: 'password reset successfully' }),
    })
  })
}

/**
 * Mock change password API.
 * mode: 'success' | 'wrong_old'
 */
export async function mockChangePasswordApi(
  page: Page,
  mode: 'success' | 'wrong_old' = 'success'
): Promise<void> {
  await page.route(`${API_BASE}/api/v1/auth/password/change`, async (route: Route) => {
    if (mode === 'success') {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'password changed' }),
      })
    } else {
      await route.fulfill({
        status: 400,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'current password is incorrect' }),
      })
    }
  })
}

/**
 * Mock user management APIs.
 */
export async function mockUserManagementApis(page: Page): Promise<void> {
  // List users
  await page.route(`${API_BASE}/api/v1/admin/users?*`, async (route: Route) => {
    const url = new URL(route.request().url())
    const search = url.searchParams.get('search')?.toLowerCase()
    const role = url.searchParams.get('role')
    const isActive = url.searchParams.get('is_active')

    let filtered = [...TEST_USERS_LIST]
    if (search) {
      filtered = filtered.filter(
        u =>
          u.first_name.toLowerCase().includes(search) ||
          u.last_name.toLowerCase().includes(search) ||
          u.email.toLowerCase().includes(search)
      )
    }
    if (role) {
      filtered = filtered.filter(u => u.role === role)
    }
    if (isActive) {
      filtered = filtered.filter(u => String(u.is_active) === isActive)
    }

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        data: {
          data: filtered,
          total: filtered.length,
          page: 1,
          limit: 20,
          total_pages: 1,
        },
      }),
    })
  })

  // Also handle without query params
  await page.route(`${API_BASE}/api/v1/admin/users`, async (route: Route) => {
    if (route.request().method() === 'GET') {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            data: TEST_USERS_LIST,
            total: TEST_USERS_LIST.length,
            page: 1,
            limit: 20,
            total_pages: 1,
          },
        }),
      })
    } else if (route.request().method() === 'POST') {
      // Create user
      const body = route.request().postDataJSON()
      const newUser = {
        id: '00000000-0000-0000-0000-000000000099',
        email: body.email,
        first_name: body.first_name,
        last_name: body.last_name,
        role: body.role,
        is_active: true,
        department_id: body.department_id ?? null,
        department_name: null,
        must_change_password: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      await route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          data: newUser,
          temporary_password: 'TempP@ss123',
        }),
      })
    } else {
      await route.fallback()
    }
  })

  // Update user
  await page.route(`${API_BASE}/api/v1/admin/users/*`, async (route: Route) => {
    if (route.request().method() === 'PUT') {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: { ...TEST_ADMIN, ...route.request().postDataJSON() } }),
      })
    } else if (route.request().method() === 'DELETE') {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'user deactivated' }),
      })
    } else {
      await route.fallback()
    }
  })

  // Admin reset password
  await page.route(`${API_BASE}/api/v1/admin/users/*/password`, async (route: Route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ message: 'password reset, user must change on next login' }),
    })
  })

  // Departments
  await page.route(`${API_BASE}/api/v1/admin/departments`, async (route: Route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ data: TEST_DEPARTMENTS }),
    })
  })
}

/**
 * Mock portal job listing APIs.
 */
export async function mockPortalJobsApi(page: Page): Promise<void> {
  // List jobs (with optional query params)
  await page.route(`${API_BASE}/api/v1/portal/jobs**`, async (route: Route) => {
    const url = new URL(route.request().url())
    const search = url.searchParams.get('search')?.toLowerCase()
    const pageNum = parseInt(url.searchParams.get('page') ?? '1')
    const limit = parseInt(url.searchParams.get('limit') ?? '12')

    let filtered = [...TEST_JOBS]
    if (search) {
      filtered = filtered.filter(
        j =>
          j.title.toLowerCase().includes(search) ||
          j.description.toLowerCase().includes(search) ||
          j.department_name.toLowerCase().includes(search)
      )
    }

    const start = (pageNum - 1) * limit
    const paged = filtered.slice(start, start + limit)

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        data: {
          data: paged,
          total: filtered.length,
          page: pageNum,
          limit,
          total_pages: Math.ceil(filtered.length / limit),
        },
      }),
    })
  })
}

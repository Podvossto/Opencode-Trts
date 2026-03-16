// E2E tests: Login flow — apps/main
// Covers: ATS-013 (Login page), ATS-003 (Login API)

import { test, expect } from '@playwright/test'
import {
  mockLoginApi,
  mockLogoutApi,
  loginAsAdmin,
  TEST_ADMIN,
} from '../fixtures/api-mocks'

test.describe('Login Page', () => {
  test.beforeEach(async ({ page }) => {
    // Clear cookies before each test
    await page.context().clearCookies()
  })

  test('renders login form with correct elements', async ({ page }) => {
    await page.goto('/login')

    // Title & description
    await expect(page.getByRole('heading', { name: 'Sign In' })).toBeVisible()
    await expect(page.getByText('Enter your credentials')).toBeVisible()

    // Form fields
    await expect(page.getByLabel('Email')).toBeVisible()
    await expect(page.getByLabel('Password')).toBeVisible()

    // Forgot password link
    await expect(page.getByRole('link', { name: 'Forgot password?' })).toBeVisible()

    // Submit button
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeVisible()
  })

  test('shows validation errors for empty fields', async ({ page }) => {
    await page.goto('/login')

    // Submit empty form
    await page.getByRole('button', { name: 'Sign In' }).click()

    // Zod validation messages should appear
    await expect(page.locator('.text-destructive').first()).toBeVisible()
  })

  test('successful login redirects to dashboard', async ({ page }) => {
    await mockLoginApi(page, 'success')
    await page.goto('/login')

    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByLabel('Password').fill('Admin123!')

    await page.getByRole('button', { name: 'Sign In' }).click()

    // Should redirect to /dashboard
    await page.waitForURL('**/dashboard**', { timeout: 10_000 })
    expect(page.url()).toContain('/dashboard')
  })

  test('invalid credentials shows error banner', async ({ page }) => {
    await mockLoginApi(page, 'invalid')
    await page.goto('/login')

    await page.getByLabel('Email').fill('wrong@email.com')
    await page.getByLabel('Password').fill('WrongPass123')

    await page.getByRole('button', { name: 'Sign In' }).click()

    // Error banner should appear
    await expect(page.getByText('Invalid email or password')).toBeVisible()
  })

  test('deactivated account shows disabled error', async ({ page }) => {
    await mockLoginApi(page, 'disabled')
    await page.goto('/login')

    await page.getByLabel('Email').fill('disabled@tigersoft.com')
    await page.getByLabel('Password').fill('SomePass123')

    await page.getByRole('button', { name: 'Sign In' }).click()

    await expect(page.getByText('Account is deactivated')).toBeVisible()
  })

  test('must_change_password user is redirected with query param', async ({ page }) => {
    await mockLoginApi(page, 'must_change')
    await page.goto('/login')

    await page.getByLabel('Email').fill('newuser@tigersoft.com')
    await page.getByLabel('Password').fill('TempPass123')

    await page.getByRole('button', { name: 'Sign In' }).click()

    // Should redirect to /dashboard?change_password=1
    await page.waitForURL('**/dashboard?change_password=1', { timeout: 10_000 })
    expect(page.url()).toContain('change_password=1')
  })

  test('session expired banner shows when reason=session_expired', async ({ page }) => {
    await page.goto('/login?reason=session_expired')

    await expect(
      page.getByText('Your session has expired. Please sign in again.')
    ).toBeVisible()
  })

  test('password visibility toggle works', async ({ page }) => {
    await page.goto('/login')

    const passwordInput = page.getByLabel('Password')
    await expect(passwordInput).toHaveAttribute('type', 'password')

    // Click eye toggle
    await page.getByRole('button', { name: 'Show password' }).click()
    await expect(passwordInput).toHaveAttribute('type', 'text')

    // Click again to hide
    await page.getByRole('button', { name: 'Hide password' }).click()
    await expect(passwordInput).toHaveAttribute('type', 'password')
  })

  test('shows loading spinner during submission', async ({ page }) => {
    // Delay the response to see spinner
    await page.route('**/api/v1/auth/login', async (route) => {
      await new Promise(r => setTimeout(r, 1000))
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            access_token: 'fake.jwt.token',
            user: TEST_ADMIN,
          },
        }),
      })
    })

    await page.goto('/login')
    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByLabel('Password').fill('Admin123!')
    await page.getByRole('button', { name: 'Sign In' }).click()

    // Spinner text should be visible
    await expect(page.getByText('Signing in...')).toBeVisible()
  })
})

test.describe('Logout', () => {
  test('unauthenticated user is redirected to login', async ({ page }) => {
    await page.context().clearCookies()
    await page.goto('/users')

    // Middleware should redirect to /login
    await page.waitForURL('**/login**', { timeout: 10_000 })
    expect(page.url()).toContain('/login')
  })
})

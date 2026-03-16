// E2E tests: Forgot Password flow — apps/main
// Covers: ATS-015 (Forgot Password 3-step OTP flow)

import { test, expect } from '@playwright/test'
import { mockOtpFlow } from '../fixtures/api-mocks'

test.describe('Forgot Password Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
    await mockOtpFlow(page)
  })

  test('renders Step 1 — email entry', async ({ page }) => {
    await page.goto('/forgot-password')

    await expect(page.getByRole('heading', { name: 'Forgot Password' })).toBeVisible()
    await expect(page.getByText("we'll send you an OTP")).toBeVisible()
    await expect(page.getByLabel('Email')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Send OTP' })).toBeVisible()
    await expect(page.getByRole('link', { name: 'Back to Sign In' })).toBeVisible()
  })

  test('Step 1 → Step 2: submitting email advances to OTP entry', async ({ page }) => {
    await page.goto('/forgot-password')

    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByRole('button', { name: 'Send OTP' }).click()

    // Step 2 should render
    await expect(page.getByRole('heading', { name: 'Enter OTP' })).toBeVisible()
    await expect(page.getByText('admin@tigersoft.com')).toBeVisible()
    await expect(page.getByLabel('Verification Code')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Verify OTP' })).toBeVisible()
  })

  test('Step 2 → Step 3: valid OTP advances to new password', async ({ page }) => {
    await page.goto('/forgot-password')

    // Step 1
    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByRole('button', { name: 'Send OTP' }).click()

    // Step 2 — enter correct OTP
    await expect(page.getByLabel('Verification Code')).toBeVisible()
    await page.getByLabel('Verification Code').fill('123456')
    await page.getByRole('button', { name: 'Verify OTP' }).click()

    // Step 3 — set new password form
    await expect(page.getByRole('heading', { name: 'Set New Password' })).toBeVisible()
    await expect(page.getByLabel('New Password')).toBeVisible()
    await expect(page.getByLabel('Confirm New Password')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Reset Password' })).toBeVisible()
  })

  test('Step 2: invalid OTP shows error', async ({ page }) => {
    await page.goto('/forgot-password')

    // Step 1
    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByRole('button', { name: 'Send OTP' }).click()

    // Step 2 — wrong OTP
    await expect(page.getByLabel('Verification Code')).toBeVisible()
    await page.getByLabel('Verification Code').fill('999999')
    await page.getByRole('button', { name: 'Verify OTP' }).click()

    // Error should appear
    await expect(page.getByText('Invalid or expired OTP')).toBeVisible()
  })

  test('Step 2: "Change email" button goes back to Step 1', async ({ page }) => {
    await page.goto('/forgot-password')

    // Step 1
    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByRole('button', { name: 'Send OTP' }).click()

    // Step 2 — click back
    await page.getByRole('button', { name: 'Change email' }).click()

    // Should be back to Step 1
    await expect(page.getByRole('heading', { name: 'Forgot Password' })).toBeVisible()
  })

  test('Step 2: "Resend code" button sends OTP again', async ({ page }) => {
    let otpRequestCount = 0
    await page.route('**/api/v1/auth/otp/request', async (route) => {
      otpRequestCount++
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'OTP sent' }),
      })
    })

    await page.goto('/forgot-password')

    // Step 1
    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByRole('button', { name: 'Send OTP' }).click()

    // Step 2 — resend
    await page.getByRole('button', { name: 'Resend code' }).click()

    // OTP request should have been called twice (initial + resend)
    expect(otpRequestCount).toBe(2)
  })

  test('Full flow: email → OTP → new password → success screen', async ({ page }) => {
    await page.goto('/forgot-password')

    // Step 1
    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByRole('button', { name: 'Send OTP' }).click()

    // Step 2
    await page.getByLabel('Verification Code').fill('123456')
    await page.getByRole('button', { name: 'Verify OTP' }).click()

    // Step 3
    await page.getByLabel('New Password').fill('NewSecure123!')
    await page.getByLabel('Confirm New Password').fill('NewSecure123!')
    await page.getByRole('button', { name: 'Reset Password' }).click()

    // Step 4 — Success screen
    await expect(page.getByRole('heading', { name: 'Password Reset Complete' })).toBeVisible()
    await expect(page.getByText('Your password has been updated')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Back to Sign In' })).toBeVisible()
  })

  test('Success screen: "Back to Sign In" navigates to /login', async ({ page }) => {
    await page.goto('/forgot-password')

    // Full flow to success
    await page.getByLabel('Email').fill('admin@tigersoft.com')
    await page.getByRole('button', { name: 'Send OTP' }).click()
    await page.getByLabel('Verification Code').fill('123456')
    await page.getByRole('button', { name: 'Verify OTP' }).click()
    await page.getByLabel('New Password').fill('NewSecure123!')
    await page.getByLabel('Confirm New Password').fill('NewSecure123!')
    await page.getByRole('button', { name: 'Reset Password' }).click()

    // Click back to sign in
    await page.getByRole('button', { name: 'Back to Sign In' }).click()
    await page.waitForURL('**/login', { timeout: 10_000 })
    expect(page.url()).toContain('/login')
  })
})

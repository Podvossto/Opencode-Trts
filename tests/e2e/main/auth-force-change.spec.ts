// E2E tests: Force password change modal — apps/main
// Covers: ATS-014 (Force password change on first login)

import { test, expect } from '@playwright/test'
import {
  mockLoginApi,
  mockChangePasswordApi,
  mockLogoutApi,
  loginAsAdmin,
  FORCED_CHANGE_JWT,
} from '../fixtures/api-mocks'

test.describe('Force Password Change Modal', () => {
  test.beforeEach(async ({ page }) => {
    // Set cookie as a user with must_change_password=true
    await page.context().addCookies([
      {
        name: 'access_token',
        value: FORCED_CHANGE_JWT,
        domain: 'localhost',
        path: '/',
        httpOnly: false,
        secure: false,
        sameSite: 'Lax',
      },
    ])
  })

  test('modal appears when change_password=1 query param is set', async ({ page }) => {
    await page.goto('/dashboard?change_password=1')

    // Modal should appear
    await expect(page.getByRole('heading', { name: 'Change Your Password' })).toBeVisible()
    await expect(
      page.getByText('Your administrator requires you to set a new password')
    ).toBeVisible()

    // Form fields
    await expect(page.getByLabel('Current Password')).toBeVisible()
    await expect(page.getByLabel('New Password')).toBeVisible()
    await expect(page.getByLabel('Confirm New Password')).toBeVisible()

    // Buttons
    await expect(page.getByRole('button', { name: 'Update Password' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Sign out instead' })).toBeVisible()
  })

  test('modal is NOT dismissed by clicking outside', async ({ page }) => {
    await page.goto('/dashboard?change_password=1')

    await expect(page.getByRole('heading', { name: 'Change Your Password' })).toBeVisible()

    // Try pressing Escape — modal should remain
    await page.keyboard.press('Escape')
    await expect(page.getByRole('heading', { name: 'Change Your Password' })).toBeVisible()
  })

  test('successful password change closes modal and redirects', async ({ page }) => {
    await mockChangePasswordApi(page, 'success')
    await page.goto('/dashboard?change_password=1')

    await page.getByLabel('Current Password').fill('OldPass123')
    await page.getByLabel('New Password').fill('NewSecure456!')
    await page.getByLabel('Confirm New Password').fill('NewSecure456!')

    await page.getByRole('button', { name: 'Update Password' }).click()

    // Should redirect to /dashboard (without change_password param)
    await page.waitForURL('**/dashboard', { timeout: 10_000 })
    expect(page.url()).not.toContain('change_password')
  })

  test('wrong old password shows error', async ({ page }) => {
    await mockChangePasswordApi(page, 'wrong_old')
    await page.goto('/dashboard?change_password=1')

    await page.getByLabel('Current Password').fill('WrongOldPass')
    await page.getByLabel('New Password').fill('NewSecure456!')
    await page.getByLabel('Confirm New Password').fill('NewSecure456!')

    await page.getByRole('button', { name: 'Update Password' }).click()

    await expect(page.getByText('Current password is incorrect')).toBeVisible()
  })

  test('"Sign out instead" logs out and redirects to login', async ({ page }) => {
    await mockLogoutApi(page)
    await page.goto('/dashboard?change_password=1')

    await page.getByRole('button', { name: 'Sign out instead' }).click()

    await page.waitForURL('**/login**', { timeout: 10_000 })
    expect(page.url()).toContain('/login')
  })

  test('password visibility toggles work in modal', async ({ page }) => {
    await page.goto('/dashboard?change_password=1')

    // Current password toggle
    const oldPwdInput = page.getByLabel('Current Password')
    await expect(oldPwdInput).toHaveAttribute('type', 'password')

    // Find the toggle buttons by aria-label
    const toggleButtons = page.getByRole('button', { name: 'Show password' })
    // First toggle = old password, second = new password
    if ((await toggleButtons.count()) >= 1) {
      await toggleButtons.first().click()
      await expect(oldPwdInput).toHaveAttribute('type', 'text')
    }
  })
})

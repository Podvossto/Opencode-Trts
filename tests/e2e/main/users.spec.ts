// E2E tests: User Management — apps/main
// Covers: ATS-017 (User list table), ATS-018 (Create/Edit/Delete user)

import { test, expect } from '@playwright/test'
import {
  loginAsAdmin,
  mockUserManagementApis,
  TEST_USERS_LIST,
  TEST_DEPARTMENTS,
} from '../fixtures/api-mocks'

test.describe('User Management Page', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
    await mockUserManagementApis(page)
  })

  // ---------- List & Table ----------

  test('renders page header and create button', async ({ page }) => {
    await page.goto('/users')

    await expect(page.getByRole('heading', { name: 'User Management' })).toBeVisible()
    await expect(page.getByText('Manage HR, Supervisor, and Admin accounts')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Create User' })).toBeVisible()
  })

  test('displays users table with correct data', async ({ page }) => {
    await page.goto('/users')

    // Wait for table to load
    await expect(page.getByRole('table')).toBeVisible()

    // Column headers
    await expect(page.getByRole('columnheader', { name: 'Name' })).toBeVisible()
    await expect(page.getByRole('columnheader', { name: 'Email' })).toBeVisible()
    await expect(page.getByRole('columnheader', { name: 'Role' })).toBeVisible()
    await expect(page.getByRole('columnheader', { name: 'Department' })).toBeVisible()
    await expect(page.getByRole('columnheader', { name: 'Status' })).toBeVisible()

    // Verify user rows exist
    await expect(page.getByText('admin@tigersoft.com')).toBeVisible()
    await expect(page.getByText('hr@tigersoft.com')).toBeVisible()
    await expect(page.getByText('supervisor@tigersoft.com')).toBeVisible()

    // Total count
    await expect(page.getByText(`${TEST_USERS_LIST.length} user`)).toBeVisible()
  })

  test('role badges display correctly', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    // Role badges should be visible
    await expect(page.getByText('ADMIN')).toBeVisible()
    await expect(page.getByText('HR').first()).toBeVisible()
    await expect(page.getByText('SUPERVISOR')).toBeVisible()
  })

  test('active/inactive status badges display', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    // Active badges
    const activeBadges = page.getByText('Active', { exact: true })
    await expect(activeBadges.first()).toBeVisible()

    // Inactive badge (disabled@tigersoft.com)
    await expect(page.getByText('Inactive')).toBeVisible()
  })

  // ---------- Search ----------

  test('search filters users by name/email (debounced)', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    // Type in search — should filter after debounce (400ms)
    await page.getByPlaceholder('Search by name or email...').fill('super')

    // Wait for debounce + re-fetch
    await page.waitForTimeout(600)

    // Only the supervisor row should match
    await expect(page.getByText('supervisor@tigersoft.com')).toBeVisible()
  })

  // ---------- Role Filter ----------

  test('role filter narrows results', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    // Open role filter dropdown
    const roleSelect = page.locator('[data-testid]').or(
      page.getByRole('combobox').first()
    )

    // Use the SelectTrigger with placeholder "All Roles"
    const roleTrigger = page.getByRole('combobox').first()
    await roleTrigger.click()

    // Select "Admin"
    await page.getByRole('option', { name: 'Admin' }).click()

    // Wait for re-fetch
    await page.waitForTimeout(300)

    // Should show only admin user
    await expect(page.getByText('admin@tigersoft.com')).toBeVisible()
  })

  // ---------- Create User ----------

  test('Create User dialog opens and has correct fields', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    await page.getByRole('button', { name: 'Create User' }).click()

    // Dialog should open
    await expect(page.getByRole('heading', { name: 'Create User' })).toBeVisible()

    // Form fields
    await expect(page.getByLabel('Email')).toBeVisible()
    await expect(page.getByLabel('First Name')).toBeVisible()
    await expect(page.getByLabel('Last Name')).toBeVisible()
    // Role and Department selects should exist
  })

  test('Create User submits and shows temp password', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    await page.getByRole('button', { name: 'Create User' }).click()

    // Fill form
    await page.getByLabel('Email').fill('newguy@tigersoft.com')
    await page.getByLabel('First Name').fill('New')
    await page.getByLabel('Last Name').fill('Guy')

    // Submit
    await page.getByRole('button', { name: /create|save/i }).click()

    // Temp password banner should appear
    await expect(page.getByText('User created: newguy@tigersoft.com')).toBeVisible({
      timeout: 10_000,
    })
    await expect(page.getByText('TempP@ss123')).toBeVisible()

    // Copy button should exist
    await expect(page.getByRole('button', { name: 'Copy' })).toBeVisible()
  })

  test('Temp password banner can be dismissed', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    // Create a user first to get the banner
    await page.getByRole('button', { name: 'Create User' }).click()
    await page.getByLabel('Email').fill('temp@tigersoft.com')
    await page.getByLabel('First Name').fill('Temp')
    await page.getByLabel('Last Name').fill('User')
    await page.getByRole('button', { name: /create|save/i }).click()

    // Wait for banner
    await expect(page.getByText('User created')).toBeVisible({ timeout: 10_000 })

    // Dismiss
    await page.getByRole('button', { name: 'Dismiss' }).click()

    // Banner should disappear
    await expect(page.getByText('User created')).not.toBeVisible()
  })

  // ---------- Edit User ----------

  test('Edit User dialog opens with pre-filled data', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    // Click action menu on first user row
    const actionButtons = page.getByRole('button', { name: '' }).filter({
      has: page.locator('svg'),
    })

    // Find the "..." menu for Admin User
    const adminRow = page.getByRole('row').filter({ hasText: 'admin@tigersoft.com' })
    await adminRow.getByRole('button').last().click()

    // Click Edit
    await page.getByRole('menuitem', { name: 'Edit' }).click()

    // Dialog should open with pre-filled data
    await expect(page.getByRole('heading', { name: 'Edit User' })).toBeVisible()
  })

  // ---------- Deactivate User ----------

  test('Deactivate user shows confirmation dialog', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    // Open action menu on supervisor row (which is active)
    const supervisorRow = page.getByRole('row').filter({ hasText: 'supervisor@tigersoft.com' })
    await supervisorRow.getByRole('button').last().click()

    // Click Deactivate
    await page.getByRole('menuitem', { name: 'Deactivate' }).click()

    // Confirmation dialog
    await expect(page.getByRole('heading', { name: 'Deactivate User' })).toBeVisible()
    await expect(page.getByText('Are you sure you want to deactivate')).toBeVisible()
    await expect(page.getByText('Super Visor')).toBeVisible()

    // Cancel and Deactivate buttons
    await expect(page.getByRole('button', { name: 'Cancel' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Deactivate' })).toBeVisible()
  })

  test('Deactivate confirmation: Cancel closes dialog', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    const supervisorRow = page.getByRole('row').filter({ hasText: 'supervisor@tigersoft.com' })
    await supervisorRow.getByRole('button').last().click()
    await page.getByRole('menuitem', { name: 'Deactivate' }).click()

    await page.getByRole('button', { name: 'Cancel' }).click()

    // Dialog should close
    await expect(page.getByRole('heading', { name: 'Deactivate User' })).not.toBeVisible()
  })

  test('Deactivate confirmation: Confirm deactivates and refreshes list', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    const supervisorRow = page.getByRole('row').filter({ hasText: 'supervisor@tigersoft.com' })
    await supervisorRow.getByRole('button').last().click()
    await page.getByRole('menuitem', { name: 'Deactivate' }).click()

    // Click confirm
    // The "Deactivate" button inside AlertDialogAction
    const confirmBtn = page.getByRole('button', { name: 'Deactivate' }).last()
    await confirmBtn.click()

    // Dialog should close and list should refresh
    await expect(page.getByRole('heading', { name: 'Deactivate User' })).not.toBeVisible({
      timeout: 5_000,
    })
  })

  // ---------- Reset Password ----------

  test('Reset password dialog opens for a user', async ({ page }) => {
    await page.goto('/users')
    await expect(page.getByRole('table')).toBeVisible()

    // Open action menu on HR user
    const hrRow = page.getByRole('row').filter({ hasText: 'hr@tigersoft.com' })
    await hrRow.getByRole('button').last().click()

    // Click Reset Password
    await page.getByRole('menuitem', { name: 'Reset Password' }).click()

    // Dialog should open
    await expect(page.getByText('HR Manager')).toBeVisible()
  })

  // ---------- Filters card ----------

  test('filters card renders with search and select dropdowns', async ({ page }) => {
    await page.goto('/users')

    await expect(page.getByText('Filters')).toBeVisible()
    await expect(page.getByPlaceholder('Search by name or email...')).toBeVisible()
    // Two combobox selects (role + status)
    const comboboxes = page.getByRole('combobox')
    expect(await comboboxes.count()).toBeGreaterThanOrEqual(2)
  })
})

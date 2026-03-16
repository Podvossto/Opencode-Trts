// E2E tests: Portal Job Listing — apps/portal
// Covers: ATS-019 (Public job listing page), ATS-012 (Public job API)

import { test, expect } from '@playwright/test'
import { mockPortalJobsApi, TEST_JOBS } from '../fixtures/api-mocks'

test.describe('Portal — Job Listing Page', () => {
  test.beforeEach(async ({ page }) => {
    await mockPortalJobsApi(page)
  })

  // ---------- Page Structure ----------

  test('renders hero section with title and description', async ({ page }) => {
    await page.goto('/jobs')

    await expect(page.getByRole('heading', { name: 'Join Our Team' })).toBeVisible()
    await expect(
      page.getByText('Find your next opportunity at TigerSoft')
    ).toBeVisible()
  })

  test('renders search bar', async ({ page }) => {
    await page.goto('/jobs')

    await expect(
      page.getByPlaceholder('Search positions by title or keyword...')
    ).toBeVisible()
  })

  test('shows total open positions count', async ({ page }) => {
    await page.goto('/jobs')

    // Wait for jobs to load
    await expect(page.getByText(`${TEST_JOBS.length} open position`)).toBeVisible({
      timeout: 10_000,
    })
  })

  // ---------- Job Cards ----------

  test('displays job cards with correct content', async ({ page }) => {
    await page.goto('/jobs')

    // Wait for cards to load
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })
    await expect(page.getByText('Frontend Engineer (React)')).toBeVisible()
    await expect(page.getByText('HR Coordinator')).toBeVisible()
    await expect(page.getByText('Data Analyst Intern')).toBeVisible()
  })

  test('job card shows department name', async ({ page }) => {
    await page.goto('/jobs')

    // Wait for content to load
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // Department names should be visible
    await expect(page.getByText('Engineering').first()).toBeVisible()
    await expect(page.getByText('Human Resources')).toBeVisible()
    await expect(page.getByText('Finance')).toBeVisible()
  })

  test('job card shows employment type and experience level', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // Meta chips
    await expect(page.getByText('Full-time').first()).toBeVisible()
    await expect(page.getByText('Internship')).toBeVisible()
    await expect(page.getByText('Senior')).toBeVisible()
    await expect(page.getByText('Mid-level')).toBeVisible()
  })

  test('job card shows slot count', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // Slots: "2 positions", "3 positions", "1 position"
    await expect(page.getByText('2 positions').first()).toBeVisible()
    await expect(page.getByText('3 positions')).toBeVisible()
    await expect(page.getByText('1 position')).toBeVisible()
  })

  test('job card has "Apply Now" button linking to apply page', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // Apply Now buttons
    const applyButtons = page.getByRole('link', { name: 'Apply Now' })
    expect(await applyButtons.count()).toBeGreaterThanOrEqual(TEST_JOBS.length)

    // First Apply Now button should link to /apply/<jobId>
    const firstLink = applyButtons.first()
    const href = await firstLink.getAttribute('href')
    expect(href).toContain('/apply/')
  })

  test('"Closing Soon" badge appears for jobs closing within 7 days', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // HR Coordinator closes on 2026-03-20 — within 7 days from 2026-03-16
    const closingSoonBadges = page.getByText('Closing Soon')
    expect(await closingSoonBadges.count()).toBeGreaterThanOrEqual(1)
  })

  // ---------- Search ----------

  test('search filters jobs by title (debounced)', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // Search for "Go"
    await page.getByPlaceholder('Search positions by title or keyword...').fill('Go')

    // Wait for debounce (400ms) + fetch
    await page.waitForTimeout(600)

    // Only "Senior Go Developer" should match
    await expect(page.getByText('Senior Go Developer')).toBeVisible()
    await expect(page.getByText('1 open position')).toBeVisible()
  })

  test('search shows matching search term in results summary', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    await page.getByPlaceholder('Search positions by title or keyword...').fill('React')
    await page.waitForTimeout(600)

    // Should show: matching "React"
    await expect(page.getByText('"React"')).toBeVisible()
  })

  // ---------- Empty State ----------

  test('empty state shows when no jobs match search', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // Search for something that doesn't exist
    await page.getByPlaceholder('Search positions by title or keyword...').fill('zzzznonexistent')
    await page.waitForTimeout(600)

    // Empty state
    await expect(page.getByText('No open positions')).toBeVisible()
    await expect(page.getByText('No positions match')).toBeVisible()

    // Clear search button
    await expect(page.getByRole('button', { name: 'Clear search' })).toBeVisible()
  })

  test('clear search button resets results', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // Search for nonexistent
    await page.getByPlaceholder('Search positions by title or keyword...').fill('zzzznonexistent')
    await page.waitForTimeout(600)

    // Click clear
    await page.getByRole('button', { name: 'Clear search' }).click()
    await page.waitForTimeout(600)

    // All jobs should be back
    await expect(page.getByText('Senior Go Developer')).toBeVisible()
    await expect(page.getByText(`${TEST_JOBS.length} open position`)).toBeVisible()
  })

  // ---------- Error State ----------

  test('error state shows retry button when API fails', async ({ page }) => {
    // Override the mock to fail
    await page.route('**/api/v1/portal/jobs**', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'internal server error' }),
      })
    })

    await page.goto('/jobs')

    // Error state
    await expect(page.getByText('Failed to load job listings')).toBeVisible({ timeout: 10_000 })
    await expect(page.getByRole('button', { name: 'Retry' })).toBeVisible()
  })

  test('retry button re-fetches jobs after error', async ({ page }) => {
    let callCount = 0

    await page.route('**/api/v1/portal/jobs**', async (route) => {
      callCount++
      if (callCount === 1) {
        // First call fails
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'internal server error' }),
        })
      } else {
        // Second call succeeds
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: {
              data: TEST_JOBS,
              total: TEST_JOBS.length,
              page: 1,
              limit: 12,
              total_pages: 1,
            },
          }),
        })
      }
    })

    await page.goto('/jobs')

    // Wait for error
    await expect(page.getByText('Failed to load job listings')).toBeVisible({ timeout: 10_000 })

    // Click retry
    await page.getByRole('button', { name: 'Retry' }).click()

    // Jobs should now load
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })
  })

  // ---------- Loading State ----------

  test('loading spinner shows while fetching', async ({ page }) => {
    // Delay the response
    await page.route('**/api/v1/portal/jobs**', async (route) => {
      await new Promise(r => setTimeout(r, 2000))
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            data: TEST_JOBS,
            total: TEST_JOBS.length,
            page: 1,
            limit: 12,
            total_pages: 1,
          },
        }),
      })
    })

    await page.goto('/jobs')

    // Loading text should appear
    await expect(page.getByText('Loading positions...')).toBeVisible()
  })

  // ---------- Grid Layout ----------

  test('job cards render in a grid layout', async ({ page }) => {
    await page.goto('/jobs')
    await expect(page.getByText('Senior Go Developer')).toBeVisible({ timeout: 10_000 })

    // Grid container should have the correct class
    const grid = page.locator('.grid')
    await expect(grid).toBeVisible()

    // Should have 4 job cards
    const cards = grid.locator(':scope > *')
    expect(await cards.count()).toBe(TEST_JOBS.length)
  })

  // ---------- Navigation ----------

  test('root path redirects to /jobs', async ({ page }) => {
    await page.goto('/')

    await page.waitForURL('**/jobs**', { timeout: 10_000 })
    expect(page.url()).toContain('/jobs')
  })
})

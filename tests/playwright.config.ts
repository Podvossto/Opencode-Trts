// Playwright config — ATS Recruitment E2E tests
// Projects: main (HR platform on :3000) + portal (Applicant portal on :3001)

import { defineConfig, devices } from '@playwright/test'

const MAIN_URL = process.env.MAIN_URL ?? 'http://localhost:3000'
const PORTAL_URL = process.env.PORTAL_URL ?? 'http://localhost:3001'
const API_URL = process.env.API_URL ?? 'http://localhost:8080'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html', { open: 'never' }],
    ['list'],
  ],
  timeout: 30_000,
  expect: {
    timeout: 10_000,
  },

  use: {
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    actionTimeout: 10_000,
  },

  projects: [
    // ---------- apps/main (HR/Admin) ----------
    {
      name: 'main',
      testDir: './e2e/main',
      use: {
        ...devices['Desktop Chrome'],
        baseURL: MAIN_URL,
      },
    },

    // ---------- apps/portal (Applicant) ----------
    {
      name: 'portal',
      testDir: './e2e/portal',
      use: {
        ...devices['Desktop Chrome'],
        baseURL: PORTAL_URL,
      },
    },
  ],

  // Start both Next.js dev servers before running tests
  webServer: [
    {
      command: 'npm run dev',
      cwd: '../apps/main',
      url: MAIN_URL,
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
    {
      command: 'npm run dev',
      cwd: '../apps/portal',
      url: PORTAL_URL,
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
  ],
})

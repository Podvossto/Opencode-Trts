// Purpose: next-intl configuration for apps/main — EN/TH locale support
// Ref: apps/main/next.config.ts (next-intl plugin is required)
// Dev must implement:
//   - supported locales: ['en', 'th'] with 'th' as default
//   - getRequestConfig() — loads locale messages from /messages/{locale}.json
//   - export routing config for use in middleware and Link components
//   - messages files must be created at: apps/main/messages/en.json, apps/main/messages/th.json

import { getRequestConfig } from 'next-intl/server'

// TODO: Dev must create messages/en.json and messages/th.json translation files

export default getRequestConfig(async ({ locale }) => ({
  messages: (await import(`../../messages/${locale}.json`)).default,
}))

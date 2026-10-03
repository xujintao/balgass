import { defineConfig } from '@playwright/test';
export default defineConfig({
  testDir: './tests/browser',
  timeout: 30000,
  fullyParallel: false,
  use: { baseURL: 'http://localhost:3000', trace: 'retain-on-failure' },
  webServer: {
    command: 'npm run dev -- --hostname localhost',
    url: 'http://localhost:3000',
    reuseExistingServer: false,
    env: {
      APP_ORIGIN: 'http://localhost:3000',
      SUPABASE_URL: process.env.SUPABASE_URL || 'http://127.0.0.1:54321',
      SUPABASE_PUBLISHABLE_KEY:
        process.env.SUPABASE_PUBLISHABLE_KEY || 'test-public-key',
      SUPABASE_SECRET_KEY:
        process.env.SUPABASE_SECRET_KEY ||
        process.env.TEST_SUPABASE_SECRET_KEY ||
        '',
      ...(process.env.E2E_LIVE_AUTH === 'true'
        ? {}
        : { NEXT_PUBLIC_TURNSTILE_SITE_KEY: 'test-site-key' }),
    },
  },
  projects: [
    { name: 'chromium', use: { browserName: 'chromium', channel: 'chromium' } },
  ],
});

import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  // Point this to your actual folder
  testDir: './tests/e2e', 
  /* Maximum time one test can run for */
  timeout: 30 * 1000,
  expect: { timeout: 5000 },
  fullyParallel: true,
  reporter: 'html',
  use: {
    /* Vite's default port */
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    video: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});

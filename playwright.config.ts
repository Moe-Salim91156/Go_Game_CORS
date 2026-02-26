import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  timeout: 60 * 1000,
  expect: { timeout: 15000 },

  // MUST be sequential — each test owns two real browser windows
  fullyParallel: false,
  workers: 1,

  reporter: 'html',

  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    // headless: false here doesn't matter — tests launch their own chromium
  },
});

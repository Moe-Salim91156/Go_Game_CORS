import { test, expect } from '@playwright/test';

test('Alice vs Bob: Automated Match', async ({ browser, request }) => {
  const roomId = `TEST_${Date.now()}`;

  // 1. Ensure Room exists in SQLite
  await request.post(`http://localhost:8080/api/game/create`, {
    data: { room_id: roomId }
  });

  const aliceCtx = await browser.newContext();
  const bobCtx = await browser.newContext();
  const alicePage = await aliceCtx.newPage();
  const bobPage = await bobCtx.newPage();

  // 2. Join via the Vite proxy
  await alicePage.goto(`/?room=${roomId}&player=1`);
  await bobPage.goto(`/?room=${roomId}&player=2`);

  // 3. FIX: Use the text visible in the snapshot ("Status: waiting") 
  // and wait for it to update to "ACTIVE"
  // We target the paragraph specifically since the class '.status-pill' was missing
  const statusLocator = alicePage.locator('p', { hasText: /Status:/i });
  
  // Increase timeout if the WebSocket handshake is slow
  await expect(statusLocator).toContainText(/active/i, { timeout: 10000 });

  // 4. Game Interaction
  const aliceBoard = alicePage.locator('button'); 
  const bobBoard = bobPage.locator('button');

  await aliceBoard.nth(0).click();
  
  // Verify sync on Bob's screen
  await expect(bobBoard.nth(0)).toHaveText('X', { timeout: 10000 });
});

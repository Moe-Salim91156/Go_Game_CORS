import { test, expect, Page, BrowserContext } from '@playwright/test';
import { chromium } from '@playwright/test';

const BASE = 'http://localhost:8080';

// ─── helpers ────────────────────────────────────────────────────────────────

async function signup(request: any, suffix: string) {
  const res = await request.post(`${BASE}/api/auth/signup`, {
    data: { username: `${suffix}_${Date.now()}`, password: 'pass123' },
  });
  const data = await res.json();
  return data.player_id as number;
}

async function createRoom(request: any, roomId: string, playerId: number) {
  await request.post(`${BASE}/api/game/create`, {
    data: { room_id: roomId, player_id: playerId },
  });
}

async function joinRoom(request: any, roomId: string, playerId: number) {
  await request.post(`${BASE}/api/game/join`, {
    data: { room_id: roomId, player_id: playerId },
  });
}

// Launch a real separate browser window at a specific screen position
async function openWindow(
  roomId: string,
  playerId: number,
  x: number,
  label: string
): Promise<{ ctx: BrowserContext; page: Page }> {
  const ctx = await chromium.launchPersistentContext('', {
    headless: false,
    args: [
      `--window-position=${x},60`,
      `--window-size=680,820`,
      `--no-first-run`,
      `--no-default-browser-check`,
      `--disable-extensions`,
      `--title=${label}`,
    ],
    viewport: { width: 660, height: 760 },
  });

  const page = ctx.pages()[0] ?? await ctx.newPage();
  await page.goto(`http://localhost:5173/?room=${roomId}&player=${playerId}`);
  return { ctx, page };
}

async function waitForStatus(page: Page, pattern: RegExp) {
  await expect(
    page.locator('p', { hasText: /Status:/i })
  ).toContainText(pattern, { timeout: 10000 });
}

async function clickCell(page: Page, index: number) {
  await page.locator('button.cell').nth(index).click();
}

async function expectCell(page: Page, index: number, value: 'X' | 'O') {
  await expect(
    page.locator('button.cell').nth(index)
  ).toHaveText(value, { timeout: 10000 });
}

// ─── shared game setup ───────────────────────────────────────────────────────

async function setupGame(request: any) {
  const roomId  = `TEST_${Date.now()}`;
  const aliceId = await signup(request, 'alice');
  const bobId   = await signup(request, 'bob');

  await createRoom(request, roomId, aliceId);

  // Alice on LEFT, Bob on RIGHT — true separate windows
  const alice = await openWindow(roomId, aliceId, 0,   '🟦 Alice (X)');
  const bob   = await openWindow(roomId, bobId,   700, '🟥 Bob   (O)');

  // Alice connects first, sees 'waiting'
  await waitForStatus(alice.page, /waiting/i);

  // Bob joins via HTTP — backend broadcasts 'active' to Alice's WS
  await joinRoom(request, roomId, bobId);

  // Both should now be active
  await waitForStatus(alice.page, /active/i);
  await waitForStatus(bob.page,   /active/i);

  return { roomId, aliceId, bobId, alice, bob };
}

async function pause(ms: number) {
  await new Promise(r => setTimeout(r, ms));
}

// ────────────────────────────────────────────────────────────────────────────
// TEST 1 — Both players connect and see active status
// ────────────────────────────────────────────────────────────────────────────
test('1. Both players connect and see active status', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  await waitForStatus(alice.page, /active/i);
  await waitForStatus(bob.page,   /active/i);

  await pause(1500); // Let you see the result
  await alice.ctx.close();
  await bob.ctx.close();
});

// ────────────────────────────────────────────────────────────────────────────
// TEST 2 — Move syncs in real-time
// ────────────────────────────────────────────────────────────────────────────
test('2. Alice move syncs to Bob in real-time', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  await clickCell(alice.page, 4);
  await expectCell(bob.page,   4, 'X');

  await pause(800);

  await clickCell(bob.page,   0);
  await expectCell(alice.page, 0, 'O');

  await pause(1500);
  await alice.ctx.close();
  await bob.ctx.close();
});

// ────────────────────────────────────────────────────────────────────────────
// TEST 3 — Turn enforcement
// ────────────────────────────────────────────────────────────────────────────
test('3. Bob cannot move on Alice turn', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  // Bob tries to go first — nothing should happen
  await clickCell(bob.page, 4);
  await pause(800);
  await expect(bob.page.locator('button.cell').nth(4)).toHaveText('');

  // Alice goes — works fine
  await clickCell(alice.page, 4);
  await expectCell(alice.page, 4, 'X');

  await pause(1500);
  await alice.ctx.close();
  await bob.ctx.close();
});

// ────────────────────────────────────────────────────────────────────────────
// TEST 4 — Occupied cell cannot be overwritten
// ────────────────────────────────────────────────────────────────────────────
test('4. Cannot click an occupied cell', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  await clickCell(alice.page, 0);
  await expectCell(bob.page, 0, 'X');
  await pause(800);

  // Bob tries to overwrite Alice's cell
  await clickCell(bob.page, 0);
  await pause(800);
  await expect(bob.page.locator('button.cell').nth(0)).toHaveText('X');

  await pause(1500);
  await alice.ctx.close();
  await bob.ctx.close();
});

// ────────────────────────────────────────────────────────────────────────────
// TEST 5 — Alice wins top row  👀
// ────────────────────────────────────────────────────────────────────────────
test('5. Alice wins with top row', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  const DELAY = 700;

  //  X | X | X
  // ---+---+---
  //  O | O |
  // ---+---+---
  //    |   |

  const moves: Array<{ page: Page; cell: number }> = [
    { page: alice.page, cell: 0 },
    { page: bob.page,   cell: 3 },
    { page: alice.page, cell: 1 },
    { page: bob.page,   cell: 4 },
    { page: alice.page, cell: 2 },
  ];

  for (const { page, cell } of moves) {
    await clickCell(page, cell);
    await pause(DELAY);
  }

  await waitForStatus(alice.page, /finished/i);
  await waitForStatus(bob.page,   /finished/i);

  await pause(2000); // Let you see the win screen
  await alice.ctx.close();
  await bob.ctx.close();
});

// ────────────────────────────────────────────────────────────────────────────
// TEST 6 — Bob wins left column  👀
// ────────────────────────────────────────────────────────────────────────────
test('6. Bob wins with left column', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  const DELAY = 700;

  //  O | X | X
  // ---+---+---
  //  O |   |
  // ---+---+---
  //  O |   |

  const moves: Array<{ page: Page; cell: number }> = [
    { page: alice.page, cell: 1 },
    { page: bob.page,   cell: 0 },
    { page: alice.page, cell: 2 },
    { page: bob.page,   cell: 3 },
    { page: alice.page, cell: 4 },
    { page: bob.page,   cell: 6 },
  ];

  for (const { page, cell } of moves) {
    await clickCell(page, cell);
    await pause(DELAY);
  }

  await waitForStatus(alice.page, /finished/i);
  await waitForStatus(bob.page,   /finished/i);

  await pause(2000);
  await alice.ctx.close();
  await bob.ctx.close();
});

// ────────────────────────────────────────────────────────────────────────────
// TEST 7 — Full draw  👀
// ────────────────────────────────────────────────────────────────────────────
test('7. Full draw game — all 9 cells filled', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  const DELAY = 600;

  // Verified draw:
  //  X | O | X
  // ---+---+---
  //  X | X | O
  // ---+---+---
  //  O | X | O

  const moves: Array<{ page: Page; cell: number }> = [
    { page: alice.page, cell: 0 },
    { page: bob.page,   cell: 1 },
    { page: alice.page, cell: 2 },
    { page: bob.page,   cell: 5 },
    { page: alice.page, cell: 3 },
    { page: bob.page,   cell: 6 },
    { page: alice.page, cell: 4 },
    { page: bob.page,   cell: 8 },
    { page: alice.page, cell: 7 },
  ];

  for (const { page, cell } of moves) {
    await clickCell(page, cell);
    await pause(DELAY);
  }

  await waitForStatus(alice.page, /finished|draw/i);
  await waitForStatus(bob.page,   /finished|draw/i);

  await pause(2000);
  await alice.ctx.close();
  await bob.ctx.close();
});

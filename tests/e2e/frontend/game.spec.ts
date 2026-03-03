import { test, expect, Page, BrowserContext } from '@playwright/test';
import { chromium } from '@playwright/test';

const API = 'http://localhost:8080';

// ─── API helpers ─────────────────────────────────────────────────────────────

async function signup(request: any, suffix: string) {
  const res = await request.post(`${API}/api/auth/signup`, {
    data: { username: `${suffix}_${Date.now()}`, password: 'pass123' },
  });
  const data = await res.json();
  return data.player_id as number;
}

async function createRoom(request: any, roomId: string, playerId: number) {
  await request.post(`${API}/api/game/create`, {
    data: { room_id: roomId, player_id: playerId },
  });
}

async function joinRoom(request: any, roomId: string, playerId: number) {
  await request.post(`${API}/api/game/join`, {
    data: { room_id: roomId, player_id: playerId },
  });
}

// ─── Browser helpers ──────────────────────────────────────────────────────────

async function openWindow(
  roomId: string,
  playerId: number,
): Promise<{ ctx: BrowserContext; page: Page }> {
  const ctx = await chromium.launchPersistentContext('', {
    headless: true,
    args: [
      `--no-first-run`,
      `--no-default-browser-check`,
      `--disable-extensions`,
    ],
    viewport: { width: 660, height: 760 },
  });

  const page = ctx.pages()[0] ?? await ctx.newPage();
  await page.goto(`http://localhost:5173/?room=${roomId}&player=${playerId}`);
  return { ctx, page };
}

function statusLocator(page: Page) {
  return page.locator('p', { hasText: /Status:/i });
}

async function waitForStatus(page: Page, pattern: RegExp) {
  await expect(statusLocator(page)).toContainText(pattern, { timeout: 10_000 });
}

function cell(page: Page, index: number) {
  return page.locator('button.cell').nth(index);
}

async function clickCell(page: Page, index: number) {
  await cell(page, index).click();
}

async function expectCell(page: Page, index: number, value: 'X' | 'O') {
  await expect(cell(page, index)).toHaveText(value, { timeout: 10_000 });
}

// ─── Shared setup ─────────────────────────────────────────────────────────────

async function setupGame(request: any) {
  const roomId  = `TEST_${Date.now()}`;
  const aliceId = await signup(request, 'alice');
  const bobId   = await signup(request, 'bob');

  await createRoom(request, roomId, aliceId);

  const alice = await openWindow(roomId, aliceId);
  const bob   = await openWindow(roomId, bobId);

  await waitForStatus(alice.page, /waiting/i);

  await joinRoom(request, roomId, bobId);

  await waitForStatus(alice.page, /active/i);
  await waitForStatus(bob.page,   /active/i);

  return { roomId, aliceId, bobId, alice, bob };
}

// ─── Tests ────────────────────────────────────────────────────────────────────

test('when both players connect they should both see the game as active', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  await waitForStatus(alice.page, /active/i);
  await waitForStatus(bob.page,   /active/i);

  await alice.ctx.close();
  await bob.ctx.close();
});

test('when alice makes a move bob should see it in real-time', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  await clickCell(alice.page, 4);
  await expectCell(bob.page, 4, 'X');

  await alice.ctx.close();
  await bob.ctx.close();
});

test('when bob makes a move alice should see it in real-time', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  await clickCell(alice.page, 4);
  await expectCell(bob.page, 4, 'X');

  await clickCell(bob.page, 0);
  await expectCell(alice.page, 0, 'O');

  await alice.ctx.close();
  await bob.ctx.close();
});

test('when bob tries to move on alices turn the cell should remain empty', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  await clickCell(bob.page, 4);
  await expect(cell(bob.page, 4)).toHaveText('', { timeout: 3_000 });

  await alice.ctx.close();
  await bob.ctx.close();
});

test('when alice marks a cell bob should not be able to overwrite it', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  await clickCell(alice.page, 0);
  await expectCell(bob.page, 0, 'X');

  await clickCell(bob.page, 0);
  await expect(cell(bob.page, 0)).toHaveText('X', { timeout: 3_000 });

  await alice.ctx.close();
  await bob.ctx.close();
});

test('when alice wins with the top row both players should see the game as finished', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  // X | X | X
  // O | O |
  const moves: Array<{ page: Page; index: number }> = [
    { page: alice.page, index: 0 },
    { page: bob.page,   index: 3 },
    { page: alice.page, index: 1 },
    { page: bob.page,   index: 4 },
    { page: alice.page, index: 2 },
  ];

  for (const { page, index } of moves) {
    await clickCell(page, index);
    // wait for the cell to reflect before the next move to avoid race conditions
    await expect(cell(page, index)).not.toHaveText('', { timeout: 5_000 });
  }

  await waitForStatus(alice.page, /finished/i);
  await waitForStatus(bob.page,   /finished/i);

  await alice.ctx.close();
  await bob.ctx.close();
});

test('when bob wins with the left column both players should see the game as finished', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  // O | X | X
  // O |   |
  // O |   |
  const moves: Array<{ page: Page; index: number }> = [
    { page: alice.page, index: 1 },
    { page: bob.page,   index: 0 },
    { page: alice.page, index: 2 },
    { page: bob.page,   index: 3 },
    { page: alice.page, index: 4 },
    { page: bob.page,   index: 6 },
  ];

  for (const { page, index } of moves) {
    await clickCell(page, index);
    await expect(cell(page, index)).not.toHaveText('', { timeout: 5_000 });
  }

  await waitForStatus(alice.page, /finished/i);
  await waitForStatus(bob.page,   /finished/i);

  await alice.ctx.close();
  await bob.ctx.close();
});

test('when all nine cells are filled with no winner both players should see a draw', async ({ request }) => {
  const { alice, bob } = await setupGame(request);

  // X | O | X
  // X | X | O
  // O | X | O
  const moves: Array<{ page: Page; index: number }> = [
    { page: alice.page, index: 0 },
    { page: bob.page,   index: 1 },
    { page: alice.page, index: 2 },
    { page: bob.page,   index: 5 },
    { page: alice.page, index: 3 },
    { page: bob.page,   index: 6 },
    { page: alice.page, index: 4 },
    { page: bob.page,   index: 8 },
    { page: alice.page, index: 7 },
  ];

  for (const { page, index } of moves) {
    await clickCell(page, index);
    await expect(cell(page, index)).not.toHaveText('', { timeout: 5_000 });
  }

  await waitForStatus(alice.page, /finished|draw/i);
  await waitForStatus(bob.page,   /finished|draw/i);

  await alice.ctx.close();
  await bob.ctx.close();
});

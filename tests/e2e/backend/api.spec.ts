import { test, expect } from '@playwright/test';

const BASE = 'http://localhost:8080';

function roomId() {
  return `TEST_${Date.now()}`;
}

async function signup(request: any, suffix: string) {
  const res = await request.post(`${BASE}/api/auth/signup`, {
    data: { username: `${suffix}_${Date.now()}`, password: 'pass123' },
  });
  expect(res.ok()).toBeTruthy();
  const data = await res.json();
  return data.player_id as number;
}

// ── Auth ─────────────────────────────────────────────────────────────────────

test('when signing up with a new username it should return a player_id', async ({ request }) => {
  const res = await request.post(`${BASE}/api/auth/signup`, {
    data: { username: `unique_${Date.now()}`, password: 'pass123' },
  });
  expect(res.status()).toBe(201);
  const body = await res.json();
  expect(body.player_id).toBeGreaterThan(0);
  expect(body.username).toContain('unique_');
});

test('when signing up with an existing username it should respond with a duplicate error', async ({ request }) => {
  const username = `dup_${Date.now()}`;
  await request.post(`${BASE}/api/auth/signup`, {
    data: { username, password: 'pass123' },
  });
  const res = await request.post(`${BASE}/api/auth/signup`, {
    data: { username, password: 'pass123' },
  });
  expect(res.status()).toBe(409);
});

// ── Room ─────────────────────────────────────────────────────────────────────

test('when a player creates a room it should return a success message with the room id', async ({ request }) => {
  const playerId = await signup(request, 'creator');
  const room = roomId();
  const res = await request.post(`${BASE}/api/game/create`, {
    data: { room_id: room, player_id: playerId },
  });
  expect(res.status()).toBe(201);
  const body = await res.json();
  expect(body.room_id).toBe(room);
});

test('when a second player joins a waiting room it should set the game status to active', async ({ request }) => {
  const aliceId = await signup(request, 'alice');
  const bobId   = await signup(request, 'bob');
  const room    = roomId();

  await request.post(`${BASE}/api/game/create`, { data: { room_id: room, player_id: aliceId } });

  const joinRes = await request.post(`${BASE}/api/game/join`, {
    data: { room_id: room, player_id: bobId },
  });
  expect(joinRes.status()).toBe(200);

  const statusRes = await request.post(`${BASE}/api/game/status`, {
    data: { room_id: room },
  });
  const body = await statusRes.json();
  expect(body.game_status).toBe('active');
  expect(body.player_x_id).toBe(aliceId);
  expect(body.player_o_id).toBe(bobId);
});

test('when fetching status of a newly created room it should be in waiting state', async ({ request }) => {
  const playerId = await signup(request, 'waiter');
  const room = roomId();

  await request.post(`${BASE}/api/game/create`, { data: { room_id: room, player_id: playerId } });

  const res = await request.post(`${BASE}/api/game/status`, { data: { room_id: room } });
  const body = await res.json();
  expect(body.game_status).toBe('waiting');
});

// ── Moves ─────────────────────────────────────────────────────────────────────

test('when a player makes a valid move it should be accepted and reflected on the board', async ({ request }) => {
  const aliceId = await signup(request, 'mover_a');
  const bobId   = await signup(request, 'mover_b');
  const room    = roomId();

  await request.post(`${BASE}/api/game/create`, { data: { room_id: room, player_id: aliceId } });
  await request.post(`${BASE}/api/game/join`,   { data: { room_id: room, player_id: bobId } });

  const moveRes = await request.post(`${BASE}/api/game/move`, {
    data: { room_id: room, player_id: aliceId, cell_index: 4 },
  });
  expect(moveRes.status()).toBe(202);

  const statusRes = await request.post(`${BASE}/api/game/status`, { data: { room_id: room } });
  const body = await statusRes.json();
  expect(body.board[4]).toBe('X');
});

test('when a player tries to move out of turn it should be rejected', async ({ request }) => {
  const aliceId = await signup(request, 'turn_a');
  const bobId   = await signup(request, 'turn_b');
  const room    = roomId();

  await request.post(`${BASE}/api/game/create`, { data: { room_id: room, player_id: aliceId } });
  await request.post(`${BASE}/api/game/join`,   { data: { room_id: room, player_id: bobId } });

  // Bob tries to go first — should fail
  const res = await request.post(`${BASE}/api/game/move`, {
    data: { room_id: room, player_id: bobId, cell_index: 0 },
  });
  expect(res.status()).toBe(400);
});

test('when a player tries to move on an occupied cell it should be rejected', async ({ request }) => {
  const aliceId = await signup(request, 'occ_a');
  const bobId   = await signup(request, 'occ_b');
  const room    = roomId();

  await request.post(`${BASE}/api/game/create`, { data: { room_id: room, player_id: aliceId } });
  await request.post(`${BASE}/api/game/join`,   { data: { room_id: room, player_id: bobId } });
  await request.post(`${BASE}/api/game/move`,   { data: { room_id: room, player_id: aliceId, cell_index: 0 } });

  const res = await request.post(`${BASE}/api/game/move`, {
    data: { room_id: room, player_id: bobId, cell_index: 0 },
  });
  expect(res.status()).toBe(400);
});

test('when alice wins with the top row the game status should be finished', async ({ request }) => {
  const aliceId = await signup(request, 'win_a');
  const bobId   = await signup(request, 'win_b');
  const room    = roomId();

  await request.post(`${BASE}/api/game/create`, { data: { room_id: room, player_id: aliceId } });
  await request.post(`${BASE}/api/game/join`,   { data: { room_id: room, player_id: bobId } });

  // X | X | X
  // O | O |
  const moves = [
    { player_id: aliceId, cell_index: 0 },
    { player_id: bobId,   cell_index: 3 },
    { player_id: aliceId, cell_index: 1 },
    { player_id: bobId,   cell_index: 4 },
    { player_id: aliceId, cell_index: 2 },
  ];
  for (const m of moves) {
    await request.post(`${BASE}/api/game/move`, { data: { room_id: room, ...m } });
  }

  const res = await request.post(`${BASE}/api/game/status`, { data: { room_id: room } });
  const body = await res.json();
  expect(body.game_status).toBe('finished');
  expect(body.winner_id).toBe(aliceId);
});

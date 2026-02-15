#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

HOST="http://localhost:8080"
ROOM_ID="TEST_$(date +%s)"

echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     TICTACTOE BACKEND TEST SUITE        ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════╝${NC}"
echo ""

# Helper function
check_response() {
    local response=$1
    local expected=$2
    local test_name=$3
    
    if echo "$response" | grep -q "$expected"; then
        echo -e "${GREEN}✓ PASS${NC} - $test_name"
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} - $test_name"
        echo "   Expected: $expected"
        echo "   Got: $response"
        return 1
    fi
}

PASS_COUNT=0
FAIL_COUNT=0

# Test 1: Signup Alice
echo -e "${YELLOW}[1/10] Testing Signup - Alice${NC}"
ALICE_RESP=$(curl -s -X POST $HOST/api/auth/signup \
    -H "Content-Type: application/json" \
    -d '{"username":"test_alice_'$RANDOM'","password":"pass123"}')

if echo "$ALICE_RESP" | grep -q "player_id"; then
    ALICE_ID=$(echo "$ALICE_RESP" | grep -o '"player_id":[0-9]*' | grep -o '[0-9]*')
    echo -e "${GREEN}✓ PASS${NC} - Alice signed up (ID: $ALICE_ID)"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Signup failed"
    echo "Response: $ALICE_RESP"
    ((FAIL_COUNT++))
    exit 1
fi
echo ""

# Test 2: Signup Bob
echo -e "${YELLOW}[2/10] Testing Signup - Bob${NC}"
BOB_RESP=$(curl -s -X POST $HOST/api/auth/signup \
    -H "Content-Type: application/json" \
    -d '{"username":"test_bob_'$RANDOM'","password":"pass456"}')

if echo "$BOB_RESP" | grep -q "player_id"; then
    BOB_ID=$(echo "$BOB_RESP" | grep -o '"player_id":[0-9]*' | grep -o '[0-9]*')
    echo -e "${GREEN}✓ PASS${NC} - Bob signed up (ID: $BOB_ID)"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Signup failed"
    echo "Response: $BOB_RESP"
    ((FAIL_COUNT++))
    exit 1
fi
echo ""

# Test 3: Create Room
echo -e "${YELLOW}[3/10] Testing Create Room${NC}"
CREATE_RESP=$(curl -s -X POST $HOST/api/game/create \
    -H "Content-Type: application/json" \
    -d "{\"room_id\":\"$ROOM_ID\",\"player_id\":$ALICE_ID}")

if check_response "$CREATE_RESP" "Room created" "Create room"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi
echo ""

# Test 4: Join Room
echo -e "${YELLOW}[4/10] Testing Join Room${NC}"
JOIN_RESP=$(curl -s -X POST $HOST/api/game/join \
    -H "Content-Type: application/json" \
    -d "{\"room_id\":\"$ROOM_ID\",\"player_id\":$BOB_ID}")

if check_response "$JOIN_RESP" "Joining Room Succesfull" "Join room"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi
echo ""

# Test 5: Check Game Status (should be active)
echo -e "${YELLOW}[5/10] Testing Game Status${NC}"
STATUS_RESP=$(curl -s -X POST $HOST/api/game/status \
    -H "Content-Type: application/json" \
    -d "{\"room_id\":\"$ROOM_ID\"}")

if check_response "$STATUS_RESP" '"game_status":"active"' "Game is active"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi
echo ""

# Test 6: Verify Player X
echo -e "${YELLOW}[6/10] Testing Player X Assignment${NC}"
if check_response "$STATUS_RESP" "\"player_x_id\":$ALICE_ID" "Alice is Player X"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi
echo ""

# Test 7: Verify Player O
echo -e "${YELLOW}[7/10] Testing Player O Assignment${NC}"
if check_response "$STATUS_RESP" "\"player_o_id\":$BOB_ID" "Bob is Player O"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi
echo ""

# Test 8: Alice makes first move
echo -e "${YELLOW}[8/10] Testing First Move (Alice - cell 4)${NC}"
MOVE_RESP=$(curl -s -X POST $HOST/api/game/move \
    -H "Content-Type: application/json" \
    -d "{\"room_id\":\"$ROOM_ID\",\"player_id\":$ALICE_ID,\"cell_index\":4}")

if check_response "$MOVE_RESP" "Move Accepted" "Alice's move"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi
echo ""

# Test 9: Verify board updated
echo -e "${YELLOW}[9/10] Testing Board Update${NC}"
STATUS_RESP=$(curl -s -X POST $HOST/api/game/status \
    -H "Content-Type: application/json" \
    -d "{\"room_id\":\"$ROOM_ID\"}")

BOARD=$(echo "$STATUS_RESP" | grep -o '"board":"[^"]*"' | cut -d'"' -f4)
if [ "${BOARD:4:1}" == "X" ]; then
    echo -e "${GREEN}✓ PASS${NC} - Board updated (X at position 4)"
    echo "   Current board: $BOARD"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Board not updated correctly"
    echo "   Board: $BOARD"
    ((FAIL_COUNT++))
fi
echo ""

# Test 10: Bob makes move
echo -e "${YELLOW}[10/10] Testing Second Move (Bob - cell 0)${NC}"
MOVE_RESP=$(curl -s -X POST $HOST/api/game/move \
    -H "Content-Type: application/json" \
    -d "{\"room_id\":\"$ROOM_ID\",\"player_id\":$BOB_ID,\"cell_index\":0}")

if check_response "$MOVE_RESP" "Move Accepted" "Bob's move"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi
echo ""

# Final board check
STATUS_RESP=$(curl -s -X POST $HOST/api/game/status \
    -H "Content-Type: application/json" \
    -d "{\"room_id\":\"$ROOM_ID\"}")

BOARD=$(echo "$STATUS_RESP" | grep -o '"board":"[^"]*"' | cut -d'"' -f4)

echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║            TEST RESULTS                  ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}Passed: $PASS_COUNT${NC}"
echo -e "${RED}Failed: $FAIL_COUNT${NC}"
echo ""
echo -e "${BLUE}Game Details:${NC}"
echo "  Room ID: $ROOM_ID"
echo "  Alice (X): $ALICE_ID"
echo "  Bob (O): $BOB_ID"
echo "  Board: $BOARD"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}╔══════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║     ✅ ALL TESTS PASSED! 🎉             ║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${YELLOW}🎮 Play in browser:${NC}"
    echo "  Alice: http://localhost:5173?room=$ROOM_ID&player=$ALICE_ID"
    echo "  Bob:   http://localhost:5173?room=$ROOM_ID&player=$BOB_ID"
    exit 0
else
    echo -e "${RED}╔══════════════════════════════════════════╗${NC}"
    echo -e "${RED}║     ❌ SOME TESTS FAILED                ║${NC}"
    echo -e "${RED}╚══════════════════════════════════════════╝${NC}"
    exit 1
fi

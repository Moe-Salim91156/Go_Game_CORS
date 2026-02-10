#!/bin/bash
HOST="http://localhost:8080"

echo "🚀 SHAZAM! Starting Multiplayer Test..."

# 1. Signup Alice (Player X)
echo "Registering Alice..."
ALICE_RESP=$(curl -s -X POST $HOST/signup \
    -H "Content-Type: application/json" \
    -d '{"username": "alice_queen", "password": "safe_password123"}')
ALICE_ID=$(echo $ALICE_RESP | jq '.player_id')

# 2. Alice Creates a Room
echo "Alice creating room SUI_777..."
curl -s -X POST $HOST/create \
    -H "Content-Type: application/json" \
    -d "{\"room_id\": \"SUI_777\", \"player_id\": $ALICE_ID}" > /dev/null

# 3. Signup Bob (Player O)
echo "Registering Bob..."
BOB_RESP=$(curl -s -X POST $HOST/signup \
    -H "Content-Type: application/json" \
    -d '{"username": "bob_builder", "password": "password456"}')
BOB_ID=$(echo $BOB_RESP | jq '.player_id')

# 4. Bob Joins Alice's Room
echo "Bob joining SUI_777..."
curl -s -X POST $HOST/join \
    -H "Content-Type: application/json" \
    -d "{\"room_id\": \"SUI_777\", \"opponent_id\": $BOB_ID}" > /dev/null

# 5. Check Live Status
echo -e "\n📊 LIVE GAME STATUS (SUI_777):"
curl -s -X POST $HOST/status \
    -H "Content-Type: application/json" \
    -d '{"room_id": "SUI_777"}' | jq .

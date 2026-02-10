#!/bin/bash
HOST="http://localhost:8080"

echo "1. Creating Room..."
curl -X POST $HOST/create -d '{"room_id": "SUI_1", "player_id": 1}'

echo -e "\n2. Joining Room..."
curl -X POST $HOST/join -d '{"room_id": "SUI_1", "opponent_id": 2}'

echo -e "\n3. Alice moves (Index 0)..."
curl -X POST $HOST/move -d '{"room_id": "SUI_1", "player_id": 1, "cell_index": 0}'

echo -e "\n4. Bob moves (Index 4)..."
curl -X POST $HOST/move -d '{"room_id": "SUI_1", "player_id": 2, "cell_index": 4}'

echo -e "\n5. Checking Status..."
curl -X POST $HOST/status -d '{"room_id": "SUI_1"}'

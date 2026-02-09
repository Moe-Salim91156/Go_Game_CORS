package main

import (
	"CorsGame/internal/models"
	"CorsGame/internal/store/sqlite"
	"fmt"
	"log"
)

func main() {
	// 1. Initialize DB
	db, err := store.OpenConnection()
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer db.Close()

	// 2. Initialize Tables
	if err := store.CreateTables(db); err != nil {
		log.Fatal("Table creation failed:", err)
	}

	// 3. Initialize Stores
	pStore := store.NewPlayerStore(db)
	gStore := store.NewGameStore(db)

	fmt.Println("--- Testing Player Store ---")

	// Test Create Player
	// Note: In your players_store.go, you pass a model.Player struct
	testUser := models.Player{Username: "Alice", Password: "password123"}
	Aliceid, err := pStore.CreatePlayer(testUser)
	fmt.Println(Aliceid)
	if err != nil {
		fmt.Println("Note: User might already exist:", err)
	} else {
		fmt.Printf("Successfully created Alice with ID: %d\n", Aliceid)
	}
	testUser2 := models.Player{Username: "bob", Password: "password123"}
	bobID, err := pStore.CreatePlayer(testUser2)
	fmt.Println(bobID)
	roomID := "TEST_TURN_FLIP"

	// 2. Create the room (Alice starts, so turn_id = 1)
	gStore.CreateNewRoom(roomID, Aliceid)
	gStore.JoinRoom(roomID, bobID) // Must join to make game 'active'

	// 3. Check INITIAL state
	initialState, _ := gStore.GetGameState(roomID)
	fmt.Printf("Before Move -> Turn ID: %d (Should be Alice/1)\n", initialState.Turn_id)

	// 4. Alice makes a move
	// This calls your SQL: SET board = ?, turn_id = CASE WHEN turn_id = player_x THEN player_o...
	err = gStore.UpdateMove(roomID, Aliceid, "X--------")
	if err != nil {
		log.Fatalf("Move failed: %v", err)
	}

	// 5. Check UPDATED state
	updatedState, _ := gStore.GetGameState(roomID)
	fmt.Printf("After Move  -> Turn ID: %d (Should be Bob/2)\n", updatedState.Turn_id)

	if updatedState.Turn_id == bobID {
		fmt.Println("SUCCESS: Turn flipped automatically in the DB!")
	} else {
		fmt.Println("FAILURE: Turn did not flip.")
	}
	err = gStore.UpdateMove(roomID, bobID, "--O-----")
	if err != nil {
		log.Fatalf("Move failed: %v", err)
	}
	updatedState, _ = gStore.GetGameState(roomID)
	fmt.Printf("After Move  -> Turn ID: %d (Should be alice/1)\n", updatedState.Turn_id)

}

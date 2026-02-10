package main

import (
	"CorsGame/internal/models"
	"CorsGame/internal/services"
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
	if err != nil {
		fmt.Println("Note: User might already exist:", err)
	} else {
		fmt.Printf("Successfully created Alice with ID: %d\n", Aliceid)
	}
	testUser2 := models.Player{Username: "bob", Password: "password123"}
	bobID, err := pStore.CreatePlayer(testUser2)
	roomID := "TEST123"

	gService := services.NewGameService(*gStore, *pStore)
	// 2. Create the room (Alice starts, so turn_id = 1)
	gService.GameStore.CreateNewRoom(roomID, Aliceid)
	gService.GameStore.JoinRoom(roomID, bobID)
	// 1. Alice moves to (0,0) - Index 0
	fmt.Println("--- Alice moves to 0 ---")
	gService.ExecuteMove(roomID, Aliceid, 0)

	// 2. Bob moves to (1,0) - Index 3
	fmt.Println("--- Bob moves to 3 ---")
	gService.ExecuteMove(roomID, bobID, 3)

	// 3. Alice moves to (0,1) - Index 1
	fmt.Println("--- Alice moves to 1 ---")
	gService.ExecuteMove(roomID, Aliceid, 1)

	// 4. Bob moves to (1,1) - Index 4
	fmt.Println("--- Bob moves to 4 ---")
	gService.ExecuteMove(roomID, bobID, 4)

	// 5. Alice moves to (0,2) - Index 2 -> WINNER!
	fmt.Println("--- Alice moves to 2 (Winning Move) ---")
	err = gService.ExecuteMove(roomID, Aliceid, 2)
	if err != nil {
		log.Fatalf("Winning move failed: %v", err)
	}

	// 6. Check the final result
	finalStatus, _ := gStore.GetGameState(roomID)
	fmt.Printf("Final Board: %s\n", finalStatus.Board)
	fmt.Printf("Game State: %s\n", finalStatus.GameState)
	fmt.Printf("Game State: %d\n", finalStatus.Winner_id)

	// // --- TEST SCENARIO: BOB WINS ---
	// roomID_BobWins := "BOB_WINS_TEST"
	// gStore.CreateNewRoom(roomID_BobWins, Aliceid)
	// gStore.JoinRoom(roomID_BobWins, bobID)

	// fmt.Println("\n--- Starting Bob Win Scenario ---")
	// gService.ExecuteMove(roomID_BobWins, Aliceid, 0) // Alice
	// gService.ExecuteMove(roomID_BobWins, bobID, 1)   // Bob
	// gService.ExecuteMove(roomID_BobWins, Aliceid, 3) // Alice
	// gService.ExecuteMove(roomID_BobWins, bobID, 4)   // Bob
	// gService.ExecuteMove(roomID_BobWins, Aliceid, 2) // Alice
	// gService.ExecuteMove(roomID_BobWins, bobID, 7)   // Bob (Winning Move!)

	// finalStatusBob, _ := gStore.GetGameState(roomID_BobWins)
	// fmt.Printf("Final Board: %s\n", finalStatusBob.Board)
	// fmt.Printf("Game State: %s\n", finalStatusBob.GameState)
	// fmt.Printf("Winner ID: %d (Should be %d)\n", finalStatusBob.Winner_id, bobID)
	// put a draw scenarior and print output to test

	// 	// --- TEST SCENARIO: DRAW ---
	// 	roomID_Draw := "DRAW_TEST"
	// 	gStore.CreateNewRoom(roomID_Draw, Aliceid)
	// 	gStore.JoinRoom(roomID_Draw, bobID)

	// 	fmt.Println("\n--- Starting Draw Scenario ---")
	// 	// Creating a "Cat's Game" pattern
	// 	moves := []int{0, 1, 2, 4, 3, 5, 7, 6, 8}
	// 	for i, cell := range moves {
	// 		pID := Aliceid
	// 		if i%2 != 0 {
	// 			pID = bobID
	// 		}
	// 		gService.ExecuteMove(roomID_Draw, pID, cell)
	// 	}

	// 	finalStatusDraw, _ := gStore.GetGameState(roomID_Draw)
	// 	fmt.Printf("Final Board: %s\n", finalStatusDraw.Board)
	// 	fmt.Printf("Game State: %s\n", finalStatusDraw.GameState)
	// 	fmt.Printf("Winner ID: %d (Should be 0/null)\n", finalStatusDraw.Winner_id)

}

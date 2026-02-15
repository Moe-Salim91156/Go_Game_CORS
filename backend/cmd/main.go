package main

import (
	"CorsGame/internal/handlers"
	"CorsGame/internal/services"
	"CorsGame/internal/store/sqlite"
	"fmt"
	"log"
	"net/http"
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
	gService := services.NewGameService(gStore, pStore)
	gHandler := handlers.NewGameHandler(*gService)

	mux := http.NewServeMux()
	mux.HandleFunc("/game/signup", gHandler.SignupHandler)
	mux.HandleFunc("/ws", gHandler.HandleWs)
	mux.HandleFunc("/game/create", gHandler.CreateRoom)
	mux.HandleFunc("/game/join", gHandler.JoinRoom)
	mux.HandleFunc("/game/move", gHandler.MoveHandler)
	mux.HandleFunc("/status", gHandler.GameStatus)

	// 4. Wrap with CORS and Start
	fmt.Println("🚀 SUI! Game Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}

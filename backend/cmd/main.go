package main

import (
	"CorsGame/internal/handlers"
	"CorsGame/internal/middleware"
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
	
	// FIXED: Grouped routes
	// Auth routes
	mux.HandleFunc("/api/auth/signup", gHandler.SignupHandler)
	
	// Game routes
	mux.HandleFunc("/api/game/create", gHandler.CreateRoom)
	mux.HandleFunc("/api/game/join", gHandler.JoinRoom)
	mux.HandleFunc("/api/game/move", gHandler.MoveHandler)
	mux.HandleFunc("/api/game/status", gHandler.GameStatus)
	
	// WebSocket
	mux.HandleFunc("/ws", gHandler.HandleWs)

	// FIXED: Apply CORS middleware
	handler := middleware.CORS(mux)

	fmt.Println("🚀 SUI! Game Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

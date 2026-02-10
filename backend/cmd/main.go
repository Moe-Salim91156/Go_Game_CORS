package main

import (
	"CorsGame/internal/handlers"
	"CorsGame/internal/services"
	"CorsGame/internal/store/sqlite"
	"fmt"
	"log"
	"net/http"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow the React frontend origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		// Handle preflight "OPTIONS" requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	gService := services.NewGameService(*gStore, *pStore)
	gHandler := handlers.NewGameHandler(*gService)

	mux := http.NewServeMux()
	mux.HandleFunc("/create", gHandler.CreateRoom)
	mux.HandleFunc("/join", gHandler.JoinRoom)
	mux.HandleFunc("/move", gHandler.MoveHandler)
	mux.HandleFunc("/status", gHandler.GameStatus)

	// 4. Wrap with CORS and Start
	fmt.Println("🚀 SUI! Game Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", CORSMiddleware(mux)))

}

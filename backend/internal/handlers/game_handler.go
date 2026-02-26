package handlers

import (
	"CorsGame/internal/services"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type GameHandler struct {
	gservice services.GameService
	GameHub  *Hub
}

func NewGameHandler(gs services.GameService) *GameHandler {
	return &GameHandler{
		gservice: gs,
		GameHub:  GameHub,
	}
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *GameHandler) HandleWs(w http.ResponseWriter, r *http.Request) {
	RoomID := r.URL.Query().Get("room")

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.GameHub.RegisterIntoRooms(RoomID, conn)
	game, err := h.gservice.GetGameState(RoomID)
	if err == nil {
		conn.WriteJSON(game)
	}
	defer h.GameHub.Unregister(RoomID, conn)
	defer conn.Close()

	for {
		var moveData struct {
			RoomID    string `json:"room_id"`
			PlayerID  int    `json:"player_id"`
			CellIndex int    `json:"cell_index"`
		}

		if err := conn.ReadJSON(&moveData); err != nil {
			log.Println("Client disconnected:", err)
			return
		}

		log.Printf("Received Move: Room=%s, Player=%d, Cell=%d", moveData.RoomID, moveData.PlayerID, moveData.CellIndex)

		err := h.gservice.ExecuteMove(moveData.RoomID, moveData.PlayerID, moveData.CellIndex)
		if err != nil {
			log.Println("Move Error ", err)
			continue
		}
		Game, _ := h.gservice.GetGameState(moveData.RoomID)
		h.GameHub.BroadCast(moveData.RoomID, Game)
	}
}

// FIXED: Using struct instead of map
func (h *GameHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID   string `json:"room_id"`
		PlayerID int    `json:"player_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := h.gservice.CreateRoom(req.RoomID, req.PlayerID)
	if err != nil {
		http.Error(w, "Could not create room", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	// FIXED: No more map[string]string
	json.NewEncoder(w).Encode(CreateRoomResponse{
		Message: "Room created succesfully",
		RoomID:  req.RoomID,
	})
}
func (h *GameHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID   string `json:"room_id"`
		PlayerID int    `json:"player_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := h.gservice.JoinRoom(req.RoomID, req.PlayerID)
	if err != nil {
		http.Error(w, "could not join room", http.StatusInternalServerError)
		return
	}

	// Broadcast updated state to anyone already connected to this room
	// This is what tells Alice's page to flip from "waiting" → "active"
	game, err := h.gservice.GetGameState(req.RoomID)
	if err == nil {
		h.GameHub.BroadCast(req.RoomID, game)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JoinRoomResponse{
		Message: "Joining Room Succesfull",
	})
}
func (h *GameHandler) MoveHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID    string `json:"room_id"`
		PlayerID  int    `json:"player_id"`
		CellIndex int    `json:"cell_index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.gservice.ExecuteMove(req.RoomID, req.PlayerID, req.CellIndex); err != nil {
		http.Error(w, "could not EXECUTE MOVE", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	// FIXED: No more map[string]string
	json.NewEncoder(w).Encode(MoveResponse{
		Message: "Move Accepted",
	})
}

func (h *GameHandler) GameStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID string `json:"room_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid Request", http.StatusBadRequest)
		return
	}
	game, err := h.gservice.GetGameByID(req.RoomID)
	if err != nil {
		http.Error(w, "could not Fetch Game Status", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	// FIXED: No more map[string]any
	json.NewEncoder(w).Encode(GameStatusResponse{
		RoomID:    req.RoomID,
		GameState: game.GameState,
		WinnerID:  game.WinnerID,
		PlayerXID: game.PlayerXID,
		PlayerOID: game.PlayerOID,
		Board:     game.Board,
	})
}

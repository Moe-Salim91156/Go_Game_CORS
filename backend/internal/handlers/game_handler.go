package handlers

import (
	"CorsGame/internal/services"
	"encoding/json"
	"net/http"
)

type GameHandler struct {
	Gs services.GameService
}

// handlers for /POST create room
// handlers for /POST rooms
// handlers for /MOVE
// join room , create room , j
func NewGameHandler(gs services.GameService) *GameHandler {
	return &GameHandler{
		Gs: gs,
	}
}

// URL/POST creat room
// create room CORS
func (h *GameHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID   string `json:"room_id"`
		PlayerID int    `json:"player_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := h.Gs.GameStore.CreateNewRoom(req.RoomID, req.PlayerID)
	if err != nil {
		http.Error(w, "Could not create room : ", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Room created succesfully",
		"room_id": req.RoomID,
	})
}

func (h *GameHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID     string `json:"room_id"`
		OpponentID int    `json:"Opponent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := h.Gs.GameStore.JoinRoom(req.RoomID, req.OpponentID)
	if err != nil {
		http.Error(w, "could not join room", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Joining Room Succesfull",
	})
}

func (h *GameHandler) MoveHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID    string `json:"Room_id"`
		PlayerID  int    `json:"Player_id"`
		CellIndex int    `json:"Cell_index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
	}
	if err := h.Gs.ExecuteMove(req.RoomID, req.PlayerID, req.CellIndex); err != nil {
		http.Error(w, "could not EXECUTE MOVE", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Move Accepted",
	})
}
func (h *GameHandler) GameStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID string `json:"Room_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid Request", http.StatusBadRequest)
	}
	game, err := h.Gs.GameStore.GetGameByID(req.RoomID)
	if err != nil {
		http.Error(w, "could not Fetch Game Status", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"RoomId":      req.RoomID,
		"Game_Status": game.GameState,
		"Winner_id":   game.Winner_id,
		"Player_x_id": game.Player_x_id,
		"Player_o_id": game.Player_o_id,
		"board":       game.Board,
	})
}

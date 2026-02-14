package handlers

import (
	"CorsGame/internal/services"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
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

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173" || origin == "http://localhost:5173/"
		// slash / stupid error
	}}

func (h *GameHandler) HandleWs(w http.ResponseWriter, r *http.Request) {
	RoomID := r.URL.Query().Get("room")

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// register the connection here ?
	GameHub.RegisterIntoRooms(RoomID, conn)
	game, err := h.Gs.GameStore.GetGameState(RoomID)
	if err == nil {
		conn.WriteJSON(game)
	}
	defer GameHub.Unregister(RoomID, conn)

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

		// DEBUG LINE FOR my Soul sanity
		log.Printf("Received Move: Room=%s, Player=%d, Cell=%d", moveData.RoomID, moveData.PlayerID, moveData.CellIndex)

		err := h.Gs.ExecuteMove(moveData.RoomID, moveData.PlayerID, moveData.CellIndex)
		if err != nil {
			log.Println("Move Error ", err)
			continue
		}
		Game, _ := h.Gs.GameStore.GetGameState(moveData.RoomID)
		GameHub.BroadCast(moveData.RoomID, Game)
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
		return
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
		return
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
		return
	}
	if err := h.Gs.ExecuteMove(req.RoomID, req.PlayerID, req.CellIndex); err != nil {
		http.Error(w, "could not EXECUTE MOVE", http.StatusBadRequest)
		return
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
		return
	}
	game, err := h.Gs.GameStore.GetGameByID(req.RoomID)
	if err != nil {
		http.Error(w, "could not Fetch Game Status", http.StatusBadRequest)
		return
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

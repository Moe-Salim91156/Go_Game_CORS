package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *GameHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	playerID, err := h.gservice.PlayerStore.CreatePlayer(req.Username, req.Password)
	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Must be FIRST
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"player_id": playerID,
		"username":  req.Username,
	})
}

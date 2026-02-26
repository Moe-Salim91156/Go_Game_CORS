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
	playerID, err := h.gservice.CreatePlayer(req.Username, req.Password)
	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// FIXED: No more map[string]any
	json.NewEncoder(w).Encode(SignupResponse{
		PlayerID: playerID,
		Username: req.Username,
	})
}

package handlers

// Response structs to replace map[string]string and map[string]any

type SignupResponse struct {
	PlayerID int    `json:"player_id"`
	Username string `json:"username"`
}

type CreateRoomResponse struct {
	Message string `json:"message"`
	RoomID  string `json:"room_id"`
}

type JoinRoomResponse struct {
	Message string `json:"message"`
}

type MoveResponse struct {
	Message string `json:"message"`
}

type GameStatusResponse struct {
	RoomID    string `json:"room_id"`
	GameState string `json:"game_status"`
	WinnerID  int    `json:"winner_id"`
	PlayerXID int    `json:"player_x_id"`
	PlayerOID int    `json:"player_o_id"`
	Board     string `json:"board"`
}

package models

type GameRoom struct {
	ID        string `json:"id"`
	Board     string `json:"board"`
	PlayerXID int    `json:"player_x_id"`
	PlayerOID int    `json:"player_o_id"`
	GameState string `json:"game_state"`
	TurnID    int    `json:"turn_id"`
	WinnerID  int    `json:"winner_id"`
}

type Player struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

package models

type GameRoom struct {
	ID          string    `json:"id"`
	Board       [9]string `json:"board"`
	Player_x_id int       `json:"player_x"`
	Player_o_id int       `json:"player_o"`
	GameState   string    `json:"game_state"`
	Turn_id     int       `json:"turn"`
	Winner_id   int       `json:"winner_id"`
}

type Player struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

package store

import (
	"CorsGame/internal/models"
	"database/sql"
	"fmt"
)

type GameStore struct {
	db *sql.DB
}

// create game room for player_x ,
// later implement joiin room where we insert palyer_o and update game state to active
func (g *GameStore) CreateNewRoom(GameId string, CreatorId string) error {
	query := `INSERT INTO GameRooms (id, player_x_id, turn_id, game_state) VALUES (?, ?, ?, 'waiting')`
	_, err := g.db.Exec(query, GameId, CreatorId, CreatorId)
	return err
}

func (g *GameStore) JoinRoom(GameId string, OpponentId int) {

}

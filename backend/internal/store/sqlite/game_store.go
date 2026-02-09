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
func (g *GameStore) CreateNewRoom(GameId string, CreatorId int) error {
	query := `INSERT INTO GameRooms (id, player_x_id, turn_id, game_state) VALUES (?, ?, ?, 'waiting')`
	_, err := g.db.Exec(query, GameId, CreatorId, CreatorId)
	return err
}

func (g *GameStore) JoinRoom(GameId string, OpponentId int) {
	// update DB add player_o_id by validating GameId
	query := `UPDATE GameRooms SET player_o_id = ?, game_state = 'active' WHERE id = ? AND game_state = 'waiting' AND player_o_id IS NULL`
	_, err := g.db.Exec(query, OpponentId, GameId)
	if err != nil {
		fmt.Printf("Error joining room: %v\n", err)
	}

}

func (s *GameStore) UpdateMove(gameID string, userID int, newBoard string) error {
	// We add 'turn_id = ?' to the WHERE clause
	query := `UPDATE GameRooms
              SET board = ?,
                  turn_id = CASE
                      WHEN turn_id = player_x_id THEN player_o_id
                      ELSE player_x_id
                  END
              WHERE id = ? AND turn_id = ? AND game_state = 'active'`

	result, err := s.db.Exec(query, newBoard, gameID, userID)
	if err != nil {
		return fmt.Errorf("failed to update move: %w", err)
	}

	// ... check RowsAffected() to see if it was actually their turn
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("not your turn or game not active")
	}
	return nil
}

func (s *GameStore) GetGameState(gameID string) (*models.GameRoom, error) {
	var gameRoom models.GameRoom
	// The order here MUST match the Scan below exactly
	query := `SELECT id, player_x_id, player_o_id, turn_id, game_state, winner_id, board 
	          FROM GameRooms WHERE id = ?`

	row := s.db.QueryRow(query, gameID)
	err := row.Scan(
		&gameRoom.ID,
		&gameRoom.Player_x_id,
		&gameRoom.Player_o_id,
		&gameRoom.Turn_id,
		&gameRoom.GameState,
		&gameRoom.Winner_id,
		&gameRoom.Board,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get game state: %w", err)
	}
	return &gameRoom, nil
}

func (g *GameStore) SetWinner(gameID string, winnerID int) error {
	// Update the game state to 'finished' and set the winner_id
	query := `UPDATE GameRooms SET game_state = 'finished', winner_id = ? WHERE id = ?`
	_, err := g.db.Exec(query, winnerID, gameID)
	return err
}

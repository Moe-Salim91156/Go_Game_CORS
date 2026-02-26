package store

import (
	"CorsGame/internal/models"
	"database/sql"
	"fmt"
)

type GameStore struct {
	db *sql.DB
}

func NewGameStore(db *sql.DB) *GameStore {
	return &GameStore{db: db}
}

func (s *GameStore) GetGameByID(id string) (*models.GameRoom, error) {
	var g models.GameRoom

	query := `SELECT id, player_x_id, player_o_id, board, turn_id, game_state FROM gameRooms WHERE id = ?`

	row := s.db.QueryRow(query, id)

	err := row.Scan(&g.ID, &g.PlayerXID, &g.PlayerOID, &g.Board, &g.TurnID, &g.GameState)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (g *GameStore) CreateNewRoom(GameId string, CreatorId int) error {
	query := `INSERT INTO GameRooms (id, player_x_id, turn_id, game_state) VALUES (?, ?, ?, 'waiting')`
	_, err := g.db.Exec(query, GameId, CreatorId, CreatorId)
	return err
}

func (g *GameStore) JoinRoom(GameId string, OpponentId int) error {
	query := `UPDATE GameRooms SET player_o_id = ?, game_state = 'active' WHERE id = ? AND game_state = 'waiting' AND player_o_id IS 0`
	_, err := g.db.Exec(query, OpponentId, GameId)
	if err != nil {
		return err
	}
	return nil

}

func (s *GameStore) UpdateMove(gameID string, userID int, newBoard string) error {
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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("not your turn or game not active")
	}
	return nil
}

// with board
func (g *GameStore) UpdateGame(gameID string, newState string, board string, winnerID int) error {
	query := `UPDATE GameRooms SET game_state = ?, board = ?, winner_id = ? WHERE id = ?`
	// ungod  this method
	_, err := g.db.Exec(query, newState, board, winnerID, gameID)
	return err
}

func (s *GameStore) GetGameState(gameID string) (*models.GameRoom, error) {
	var gameRoom models.GameRoom
	query := `SELECT id, player_x_id, player_o_id, turn_id, game_state, winner_id, board 
	          FROM GameRooms WHERE id = ?`

	row := s.db.QueryRow(query, gameID)
	err := row.Scan(
		&gameRoom.ID,
		&gameRoom.PlayerXID,
		&gameRoom.PlayerOID,
		&gameRoom.TurnID,
		&gameRoom.GameState,
		&gameRoom.WinnerID,
		&gameRoom.Board,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get game state: %w", err)
	}
	return &gameRoom, nil
}

func (g *GameStore) SetWinner(gameID string, winnerID int) error {
	query := `UPDATE GameRooms SET game_state = 'finished', winner_id = ? WHERE id = ?`
	_, err := g.db.Exec(query, winnerID, gameID)
	return err
}

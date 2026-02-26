package store

import (
	"CorsGame/internal/models"
	"database/sql"
	"fmt"
)

type PlayerStore struct {
	db *sql.DB
}

func NewPlayerStore(db *sql.DB) *PlayerStore {
	return &PlayerStore{db: db}
}

// Fixed: Removed HashPassword (moved to service)
// Now accepts already-hashed password
func (p *PlayerStore) CreatePlayer(Username string, HashedPassword string) (int, error) {
	query := `INSERT INTO Players (username , password) VALUES (?, ?)`

	result, err := p.db.Exec(query, Username, HashedPassword)
	if err != nil {
		return 0, fmt.Errorf("Could not execute the create player query")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (p *PlayerStore) GetPlayerByUsername(username string) (models.Player, error) {
	var player models.Player
	query := `SELECT id , username , password FROM players WHERE username = ?`
	row := p.db.QueryRow(query, username)

	err := row.Scan(&player.ID, &player.Username, &player.Password)
	if err != nil {
		return player, err
	}
	return player, nil
}

func (p *PlayerStore) GetPlayerById(id int) (models.Player, error) {
	var player models.Player
	query := `SELECT id , username , password FROM players WHERE id = ?`
	row := p.db.QueryRow(query, id)

	err := row.Scan(&player.ID, &player.Username, &player.Password)
	if err != nil {
		return player, err
	}
	return player, nil
}

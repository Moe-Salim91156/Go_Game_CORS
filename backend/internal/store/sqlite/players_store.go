package store

import (
	"CorsGame/internal/models"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type PlayerStore struct {
	db *sql.DB
}

func NewPlayerStore(db *sql.DB) *PlayerStore {
	return &PlayerStore{db: db}
}

// create User

// HashPassword converts a plain text password into a Bcrypt hash.
func HashPassword(password string) (string, error) {
	// Cost of 10-12 is standard for 2026. Higher = slower but more secure.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// for now its here , unitl i implelent auth services

// // CheckPasswordHash compares a plain text password with a stored hash.
// func CheckPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

func (p *PlayerStore) CreatePlayer(player models.Player) (int, error) {

	query := `INSERT INTO Players (username , password) VALUES (?, ?)`

	hashed_password, err := HashPassword(player.Password)
	result, err := p.db.Exec(query, player.Username, hashed_password)
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
	// input uername string , select from db and return
	var player models.Player
	query := `SELECT (id , username , password FROM players WHERE username = ?)`
	row := p.db.QueryRow(query, username)

	err := row.Scan(&player.ID, &player.Username, &player.Password)
	if err != nil {
		return player, err
	}
	return player, nil
}

func (p *PlayerStore) GetPlayerById(id int) (models.Player, error) {
	// input id integer, select from db and return
	var player models.Player
	query := `SELECT (id , username , password FROM players WHERE id = ?)`
	row := p.db.QueryRow(query, id)

	err := row.Scan(&player.ID, &player.Username, &player.Password)
	if err != nil {
		return player, err
	}
	return player, nil
}

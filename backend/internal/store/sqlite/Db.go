package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// init sqlite
func OpenConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// -> init tables
func CreateTables(db *sql.DB) error {

	userTableQuery := `CREATE TABLE IF NOT EXISTS Players (
	id  INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);`

	gameTableQuery := `CREATE TABLE IF NOT EXISTS GameRooms (
id TEXT PRIMARY KEY,
player_x_id INTEGER NOT NULL,
player_o_id INTEGER DEFAULT 0,
board TEXT DEFAULT '---------',
turn_id INTEGER NOT NULL,
game_state TEXT DEFAULT 'waiting',
winner_id INTEGER DEFAULT 0,
FOREIGN KEY(player_x_id) REFERENCES players(id),
FOREIGN KEY(player_o_id) REFERENCES players(id),
FOREIGN KEY(turn_id) REFERENCES players(id)
);`

	_, err := db.Exec(userTableQuery)
	if err != nil {
		return err
	}
	_, err = db.Exec(gameTableQuery)
	if err != nil {
		return err
	}
	return nil
}

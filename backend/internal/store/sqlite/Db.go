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

	UserTableQuery := `CREATE TABLE IF NOT EXISTS Players (
	id  INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);`

	GameTableQuery := `CREATE TABLE IF NOT EXISTS GameRooms (
id TEXT PRIMARY KEY,
player_x_id INTEGER,
player_o_id INTEGER,
board TEXT DEFAULT '---------',
turn_id INTEGER,
game_state TEXT DEFAULT 'waiting',
winner_id INTEGER,
FOREIGN KEY(player_x_id) REFERENCES players(id),
FOREIGN KEY(player_o_id) REFERENCES players(id),
FOREIGN KEY(turn_id) REFERENCES players(id)
);`

	_, err := db.Exec(UserTableQuery)
	if err != nil {
		return err
	}
	_, err = db.Exec(GameTableQuery)
	if err != nil {
		return err
	}
	return nil
}

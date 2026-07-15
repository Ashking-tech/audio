package storage

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func InitializeDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`

		CREATE TABLE IF NOT EXISTS songs(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);
		CREATE TABLE IF NOT EXISTS fingerprints(
			hash INTEGER,
			anchor_time INTEGER,
			song_id INTEGER,
			FOREIGN KEY (song_id) REFERENCES songs(id)
		);
		CREATE INDEX IF NOT EXISTS idx_hash ON fingerprints(hash);
	`)
	if err != nil {
		return nil, err
	}
	return db, nil

}

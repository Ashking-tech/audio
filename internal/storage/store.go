package storage

import (
	"database/sql"
	"fmt"

	"github.com/Ashking-tech/audio/internal/fingerprint"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func Insertsong(database *sql.DB, name string, fps []fingerprint.Fingerprint) (int64, error) {

	result, err := database.Exec("insert into songs (name) values (?)", name)
	if err != nil {
		return 0, fmt.Errorf("insert songs: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get songs id : %w", err)
	}

	return id, InsertFingerprints(database, id, fps)

}

func InsertFingerprints(database *sql.DB, songID int64, fps []fingerprint.Fingerprint) error {
	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction : %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO fingerprints(hash,anchor_time, song_id) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare statement %w", err)
	}
	defer stmt.Close()

	for _, fp := range fps {
		_, err := stmt.Exec(fp.Hash, fp.AnchorTime, songID)
		if err != nil {
			return fmt.Errorf("insert fingerprint :%w", err)
		}
	}
	return tx.Commit()
}

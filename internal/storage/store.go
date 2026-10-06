package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Ashking-tech/audio/internal/fingerprint"
	// Turso rewrite: drivers are registered in db.go (sqlite + libsql).
	// This file needs no driver import to function, but we keep the blank
	// imports so it compiles standalone in tests/tools.
	_ "github.com/tursodatabase/go-libsql"
	_ "modernc.org/sqlite"
)

var db *sql.DB

// ---------------------------------------------------------------------------
// ORIGINAL (pre-Turso) implementations — kept for reference, commented out.
//
// func Insertsong(database *sql.DB, name string, fps []fingerprint.Fingerprint) (int64, error) {
//
// 	result, err := database.Exec("insert into songs (name) values (?)", name)
// 	if err != nil {
// 		return 0, fmt.Errorf("insert songs: %w", err)
// 	}
//
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("get songs id : %w", err)
// 	}
//
// 	return id, InsertFingerprints(database, id, fps)
//
// }
//
// func InsertFingerprints(database *sql.DB, songID int64, fps []fingerprint.Fingerprint) error {
// 	tx, err := database.Begin()
// 	if err != nil {
// 		return fmt.Errorf("begin transaction : %w", err)
// 	}
// 	defer tx.Rollback()
//
// 	stmt, err := tx.Prepare("INSERT INTO fingerprints(hash,anchor_time, song_id) VALUES (?, ?, ?)")
// 	if err != nil {
// 		return fmt.Errorf("prepare statement %w", err)
// 	}
// 	defer stmt.Close()
//
// 	for _, fp := range fps {
// 		_, err := stmt.Exec(fp.Hash, fp.AnchorTime, songID)
// 		if err != nil {
// 			return fmt.Errorf("insert fingerprint :%w", err)
// 		}
// 	}
// 	return tx.Commit()
// }
// ---------------------------------------------------------------------------

// Insertsong inserts the song row and delegates fingerprint writes.
// Turso rewrite note: UNCHANGED logic. Both sqlite and libsql support
// LastInsertId() (backed by last_insert_rowid()), so this works on local
// and Turso without branching.
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

// InsertFingerprints bulk-inserts fingerprints in multi-row chunks.
//
// Turso rewrite notes (WHY this changed):
//   - Original: 1 prepared statement + N x Exec (N ~= 40k per song).
//     Local SQLite: fast (same process, ~1s).
//     Remote Turso: 40k network round-trips -> minutes per add-song,
//     plus 10M monthly row-write quota burn from chatty protocol.
//   - New: chunk into multi-row INSERTs of 500 rows each:
//     "INSERT ... VALUES (?,?,?),(?,?,?)... x500" -> ~80 round-trips/song.
//     Same rows, same transaction atomicity, ~500x fewer trips.
//   - Chunk size 500: safely under libsql/SQLite variable limits
//     (500 rows x 3 vars = 1500 vars; libsql allows more, but 500 keeps
//     each SQL statement small and retry-friendly).
//   - Empty input: no-op, returns nil (original would commit empty tx).
func InsertFingerprints(database *sql.DB, songID int64, fps []fingerprint.Fingerprint) error {
	if len(fps) == 0 {
		return nil
	}

	const chunkSize = 500

	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction : %w", err)
	}
	defer tx.Rollback()

	for start := 0; start < len(fps); start += chunkSize {
		end := start + chunkSize
		if end > len(fps) {
			end = len(fps)
		}
		chunk := fps[start:end]

		var sb strings.Builder
		sb.WriteString("INSERT INTO fingerprints(hash,anchor_time, song_id) VALUES ")
		args := make([]any, 0, len(chunk)*3)
		for i, fp := range chunk {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString("(?, ?, ?)")
			args = append(args, fp.Hash, fp.AnchorTime, songID)
		}

		if _, err := tx.Exec(sb.String(), args...); err != nil {
			return fmt.Errorf("insert fingerprint chunk [%d:%d]: %w", start, end, err)
		}
	}
	return tx.Commit()
}

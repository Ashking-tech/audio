package storage

import (
	"database/sql"
	"os"
	"strings"

	// Turso rewrite: go-libsql registers the "libsql" database/sql driver.
	// It speaks the same database/sql API as modernc.org/sqlite, so the rest
	// of the code (Exec/Query/Begin) stays unchanged. We keep both drivers:
	//   - "sqlite"  -> local file (CLI use, tests with ":memory:")
	//   - "libsql"  -> remote Turso (prod, when TURSO_URL is set)
	_ "github.com/tursodatabase/go-libsql"
	_ "modernc.org/sqlite"
)

// ---------------------------------------------------------------------------
// ORIGINAL (pre-Turso) implementation — kept for reference, commented out.
//
// func InitializeDB(path string) (*sql.DB, error) {
// 	db, err := sql.Open("sqlite", path)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	_, err = db.Exec(`
//
// 		CREATE TABLE IF NOT EXISTS songs(
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			name TEXT UNIQUE
// 		);
// 		CREATE TABLE IF NOT EXISTS fingerprints(
// 			hash INTEGER,
// 			anchor_time INTEGER,
// 			song_id INTEGER,
// 			FOREIGN KEY (song_id) REFERENCES songs(id)
// 		);
// 		CREATE INDEX IF NOT EXISTS idx_hash ON fingerprints(hash);
// 	`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
//
// }
// ---------------------------------------------------------------------------

// InitializeDB opens either Turso (remote) or local SQLite.
//
// Turso rewrite notes:
//   - Local behaviour is UNCHANGED: InitializeDB("fingerprints.db") and
//     InitializeDB(":memory:") (used by tests) still use the "sqlite" driver.
//   - Prod behaviour: if TURSO_URL is set (e.g. "libsql://xxx.turso.io"),
//     we use the "libsql" driver instead. TURSO_TOKEN is appended as
//     ?authToken=... unless the URL already carries it.
//   - Schema is identical on both backends (songs + fingerprints + idx_hash),
//     so local and remote DBs are interchangeable.
//   - Env vars (set these on Render/Koyeb/Fly, NOT in git):
//     TURSO_URL   = libsql://<db>-<org>.turso.io
//     TURSO_TOKEN = libsql auth token (database scoped recommended)
//     DB_PATH     = handled by caller (main.go), passed in as `path`.
func InitializeDB(path string) (*sql.DB, error) {
	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_TOKEN")

	var db *sql.DB
	var err error

	if tursoURL != "" {
		// Remote Turso path.
		connStr := tursoURL
		if tursoToken != "" && !strings.Contains(connStr, "authToken") {
			sep := "?"
			if strings.Contains(connStr, "?") {
				sep = "&"
			}
			connStr = connStr + sep + "authToken=" + tursoToken
		}
		db, err = sql.Open("libsql", connStr)
	} else {
		// Local SQLite path (original behaviour).
		if path == "" {
			path = "fingerprints.db"
		}
		db, err = sql.Open("sqlite", path)
	}
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

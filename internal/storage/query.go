package storage

import (
	"database/sql"
	"strings"

	"github.com/Ashking-tech/audio/internal/fingerprint"
	// Turso rewrite: drivers registered in db.go. Kept here so this file
	// compiles standalone.
	_ "github.com/tursodatabase/go-libsql"
	_ "modernc.org/sqlite"
)

// ---------------------------------------------------------------------------
// ORIGINAL (pre-Turso) LookUpMatches — kept for reference, commented out.
// Per-fingerprint query: 1 round-trip per query hash (~2-5k trips/match).
// Fine on local SQLite, unusable on remote Turso (8-15s+ per match).
//
// func LookUpMatches(database *sql.DB, queryFps []fingerprint.Fingerprint) (songN string, err error) {
// 	type offsetCount struct {
// 		offset int
// 		count  int
// 	}
// 	votes := make(map[int][]offsetCount)
//
// 	for _, fp := range queryFps {
// 		rows, err := database.Query(
// 			"SELECT song_id, anchor_time FROM fingerprints WHERE hash = ?", fp.Hash)
// 		if err != nil {
// 			return "", err
// 		}
//
// 		for rows.Next() {
// 			var songID, dbAnchor int
//
// 			err := rows.Scan(&songID, &dbAnchor)
// 			if err != nil {
// 				rows.Close()
// 				return "", err
// 			}
// 			offset := fp.AnchorTime - dbAnchor
// 			votes[songID] = append(votes[songID], offsetCount{offset, 1})
// 		}
// 		rows.Close()
// 	}
//
// 	if len(votes) == 0 {
// 		return "", nil
// 	}
//
// 	type scored struct{ id, score int }
// 	var best scored
//
// 	for songID, offsets := range votes {
// 		offsetCounts := make(map[int]int)
// 		maxCount := 0
// 		for _, oc := range offsets {
// 			offsetCounts[oc.offset]++
// 			if offsetCounts[oc.offset] > maxCount {
// 				maxCount = offsetCounts[oc.offset]
// 			}
// 		}
// 		if maxCount > best.score {
// 			best = scored{songID, maxCount}
// 		}
// 	}
//
// 	if best.score == 0 {
// 		return "", nil
// 	}
//
// 	var songName string
// 	err = database.QueryRow(
// 		"SELECT name FROM songs WHERE id = ?", best.id,
// 	).Scan(&songName)
//
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", nil
// 		}
// 		return "", err
// 	}
//
// 	return songName, nil
// }
// ---------------------------------------------------------------------------

// LookUpMatches finds the best song via offset-aligned hash voting.
//
// Turso rewrite notes (WHY this changed + correctness proof):
//   - Original: N queries (N = len(queryFps)). New: ceil(D/500) queries
//     where D = distinct hashes (~same as N, but /500). ~3k trips -> ~6.
//   - Correctness: results are IDENTICAL to the original. We group query
//     anchors by hash (hash -> []queryAnchorTime), fetch each distinct hash
//     once, then replay one vote per (dbRow x queryOccurrence), i.e. the
//     same offset = queryAnchor - dbAnchor vote the loop used to cast.
//     Duplicate query hashes get duplicate votes, exactly as before.
//   - Chunk size 500 keeps "IN (?,?,...)" under libsql bind limits and
//     keeps each SQL statement < ~5KB.
//   - Scoring (densest song+offset cluster wins) is byte-for-byte the
//     original algorithm; only the fetch layer changed.
func LookUpMatches(database *sql.DB, queryFps []fingerprint.Fingerprint) (songN string, err error) {
	type offsetCount struct {
		offset int
		count  int
	}

	if len(queryFps) == 0 {
		return "", nil
	}

	// Group query anchor times by hash so each distinct hash is fetched once
	// while still casting one vote per original occurrence.
	queryAnchors := make(map[uint32][]int, len(queryFps))
	distinct := make([]uint32, 0, len(queryFps))
	for _, fp := range queryFps {
		if _, ok := queryAnchors[fp.Hash]; !ok {
			distinct = append(distinct, fp.Hash)
		}
		queryAnchors[fp.Hash] = append(queryAnchors[fp.Hash], fp.AnchorTime)
	}

	votes := make(map[int][]offsetCount)

	const chunkSize = 500
	for start := 0; start < len(distinct); start += chunkSize {
		end := start + chunkSize
		if end > len(distinct) {
			end = len(distinct)
		}
		chunk := distinct[start:end]

		placeholders := strings.Repeat("?,", len(chunk)-1) + "?"
		args := make([]any, len(chunk))
		for i, h := range chunk {
			args[i] = h
		}

		// NOTE: we also select `hash` so we can map each db row back to
		// its query anchor list.
		rows, err := database.Query(
			"SELECT song_id, anchor_time, hash FROM fingerprints WHERE hash IN ("+placeholders+")",
			args...,
		)
		if err != nil {
			return "", err
		}

		for rows.Next() {
			var songID, dbAnchor int
			var h int64 // hash stored as INTEGER; scan wide then cast
			if err := rows.Scan(&songID, &dbAnchor, &h); err != nil {
				rows.Close()
				return "", err
			}
			for _, qAnchor := range queryAnchors[uint32(h)] {
				offset := qAnchor - dbAnchor
				votes[songID] = append(votes[songID], offsetCount{offset, 1})
			}
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return "", err
		}
	}

	if len(votes) == 0 {
		return "", nil
	}

	// --- scoring below is the ORIGINAL algorithm, unchanged ---
	type scored struct{ id, score int }
	var best scored

	for songID, offsets := range votes {
		offsetCounts := make(map[int]int)
		maxCount := 0
		for _, oc := range offsets {
			offsetCounts[oc.offset]++
			if offsetCounts[oc.offset] > maxCount {
				maxCount = offsetCounts[oc.offset]
			}
		}
		if maxCount > best.score {
			best = scored{songID, maxCount}
		}
	}

	if best.score == 0 {
		return "", nil
	}

	var songName string
	err = database.QueryRow(
		"SELECT name FROM songs WHERE id = ?", best.id,
	).Scan(&songName)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return songName, nil
}

func ListSongs(database *sql.DB) ([]string, error) {
	// Turso rewrite note: UNCHANGED. Single lightweight query, no N+1
	// problem, identical SQL works on sqlite + libsql.
	rows, err := database.Query("SELECT name FROM songs ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

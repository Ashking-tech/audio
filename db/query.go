package db

import (
	"database/sql"

	"github.com/Ashking-tech/audio/fingerprint"
	_ "modernc.org/sqlite"
)

func LookUpMatches(db *sql.DB, queryFps []fingerprint.Fingerprint) (songN string, err error) {
	type offsetCount struct {
		offset int
		count  int
	}
	votes := make(map[int][]offsetCount)

	for _, fp := range queryFps {
		rows, err := db.Query(
			"SELECT song_id, anchor_time FROM fingerprints WHERE hash = ?", fp.Hash)
		if err != nil {
			return "", err
		}

		for rows.Next() {
			var songID, dbAnchor int

			err := rows.Scan(&songID, &dbAnchor)
			if err != nil {
				rows.Close()
				return "", err
			}
			offset := fp.AnchorTime - dbAnchor
			votes[songID] = append(votes[songID], offsetCount{offset, 1})
		}
		rows.Close()
	}

	if len(votes) == 0 {
		return "", nil
	}

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
	err = db.QueryRow(
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
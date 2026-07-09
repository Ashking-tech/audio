package db

import (
	"database/sql"

	"github.com/Ashking-tech/audio/fingerprint"
	_ "modernc.org/sqlite"
)

func LookUpMatches( db *sql.DB, queryFps []fingerprint.Fingerprint) (songN string ,err error) {
		voteCount := make(map[int]int)
	for _,fp := range queryFps {
		result,err := db.Query(
			"SELECT song_id FROM Fingerprints WHERE hash = ?",fp.Hash)
		if err != nil {
			return "",err
		}

		for result.Next(){
			var songID int 

			err := result.Scan(&songID)
			if err != nil {
				return "",err
			}
			voteCount[songID]++
			
		}

		result.Close()
	}

	if len(voteCount) == 0 {
		return "", nil
	}
	bestSongID := 0
	bestVotes := 0

 for songID,votes := range voteCount {
		if votes > bestVotes {
			bestVotes = votes
			bestSongID = songID
		}
	}
	var songName string

err = db.QueryRow(
		"SELECT name FROM songs WHERE id = ?",
		bestSongID,
).Scan(&songName)

if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
}

return songName, nil
	
}
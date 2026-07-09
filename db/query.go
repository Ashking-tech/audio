package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
	"github.com/Ashking-tech/audio/fingerprint"
)

func LookUpMatches( db *sql.DB, queryFps []fingerprint.Fingerprint) (string songN,err error) {
		voteCount := make(map[int]int)
		var songID int
	for _,fp := range queryFps {
		result,err := db.Query(
			"SELECT song_id FROM Fingerprints WHERE hash = ?",fp.Hash)
		if err != nil {
			return "",err
		}
		defer result.Close()

		for result.Next()
	}
	
}
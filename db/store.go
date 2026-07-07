package db


import (
	"database/sql"
	_"modernc.org/sqlite"
)

var db *sql.DB

func InsertSong(db,name,fps)(){
	_,err := db.Exec(

		"INSERT INTO fingerprints(hash,song_id, anchor_time) VALUES (?,?,?)",
		hash,
		songID,
		anchorTime,
	)
return err
}
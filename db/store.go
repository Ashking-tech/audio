package db

import (
	"database/sql"
	"fmt"

	"github.com/Ashking-tech/audio/fingerprint"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func insertsong(db,name,fps)(){

	result, err := db.Exec("insert into songs (name) values (?)", name)
	if err != nil {
		return 0, fmt.Errorf("insert songs: %w",err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get songs id : %w",err)
	}

	return id,nil

	
	
}

func InsertFingerprints(db *sql.DB,songID int64,fps[] fingerprint.Fingerprint) error {
	tx,err := db.Begin() //get ready nigger
	if err != nil {
		return fmt.Errorf("begin transaction : %w",err)
	}
	defer tx.Rollback() //does nothing if tx.commit already called

	//put the sql query inside the stmt and then execute it inside the for loop,in that execution pass the values like function parameters
	stmt,err := tx.Prepare("INSERT INTO fingerprints(hash,anchor_time, song_id) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare statement %w",err)
	}
	defer stmt.Close()

	
	for _,fp := range fps { //this for loops does the job of insertion rahhhhhhhh
		_,err := stmt.Exec(fp.Hash,fp.AnchorTime,songID) //right here the values are inserted inside the stmt query
		if err != nil {
			return fmt.Errorf("insert fingerprint :%w",err)
		}
	}
	return tx.Commit()
}
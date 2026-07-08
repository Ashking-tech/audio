package db

import (
	"database/sql"
	"fmt"

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
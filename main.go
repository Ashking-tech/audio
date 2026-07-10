package main

import (
	"fmt"

	"github.com/Ashking-tech/audio/db"
	"github.com/Ashking-tech/audio/pipeline"
)

func main() {

	database,_ := db.InitializeDB("fingerprints.db")
	defer database.Close()
	// err := pipeline.IngestPipeline(database,"output.wav","output.wav")
	// if err != nil {
	// 	panic(err)
	// }

	match, err := pipeline.MatchFile(database,"output.wav")
	if err != nil {
		panic(err)
	}
	fmt.Println("best match: ",match)
}

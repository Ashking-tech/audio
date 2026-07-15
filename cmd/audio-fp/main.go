package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Ashking-tech/audio/internal/audio"
	"github.com/Ashking-tech/audio/internal/pipeline"
	"github.com/Ashking-tech/audio/internal/storage"
	"github.com/Ashking-tech/audio/internal/web"
)

func main() {
	database, err := storage.InitializeDB("fingerprints.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  add-song <file.wav> <name>   — add a song to the database")
		fmt.Println("  match <file.wav>             — match a file against the database")
		fmt.Println("  listen [duration]            — record from mic and match")
		fmt.Println("  list                         — list all songs")
		fmt.Println("  serve                        — start the web server")
		return
	}

	switch os.Args[1] {
	case "add-song":
		if len(os.Args) < 4 {
			log.Fatal("Usage: add-song <file.wav> <name>")
		}
		err := pipeline.IngestPipeline(database, os.Args[2], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("added:", os.Args[3])

	case "match":
		if len(os.Args) < 3 {
			log.Fatal("Usage: match <file.wav>")
		}
		match, err := pipeline.MatchFile(database, os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		if match == "" {
			fmt.Println("no match found")
		} else {
			fmt.Println("best match:", match)
		}

	case "listen":
		duration := 5
		if len(os.Args) >= 3 {
			duration, err = strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatal("invalid duration")
			}
		}
		fmt.Println("recording for", duration, "seconds...")
		samples, err := audio.Record(duration, 44100)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("recorded", len(samples), "samples")
		match, err := pipeline.MatchRecording(database, samples)
		if err != nil {
			log.Fatal(err)
		}
		if match == "" {
			fmt.Println("no match found")
		} else {
			fmt.Println("best match:", match)
		}

	case "list":
		songs, err := storage.ListSongs(database)
		if err != nil {
			log.Fatal(err)
		}
		if len(songs) == 0 {
			fmt.Println("no songs in database")
		} else {
			for i, name := range songs {
				fmt.Printf("%d. %s\n", i+1, name)
			}
		}

	case "serve":
		srv := web.NewServer(database, ":8082")
		log.Fatal(srv.Start())

	default:
		fmt.Println("unknown command:", os.Args[1])
	}
}

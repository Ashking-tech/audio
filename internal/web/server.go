package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Ashking-tech/audio/internal/pipeline"
	"github.com/Ashking-tech/audio/internal/storage"
)

type Server struct {
	db   *sql.DB
	addr string
	tmp  string
}

func NewServer(database *sql.DB, addr string) *Server {
	tmp := filepath.Join(os.TempDir(), "audio-fingerprint")
	os.MkdirAll(tmp, 0755)
	return &Server{db: database, addr: addr, tmp: tmp}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web/static")))
	mux.HandleFunc("/api/add-song", s.apiAddSong)
	mux.HandleFunc("/api/match", s.apiMatch)
	mux.HandleFunc("/api/songs", s.apiSongs)
	log.Println("listening on", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) saveUpload(r *http.Request) (string, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("form file: %w", err)
	}
	defer file.Close()

	tmp, err := os.CreateTemp(s.tmp, "upload-*.wav")
	if err != nil {
		return "", fmt.Errorf("temp file: %w", err)
	}

	_, err = io.Copy(tmp, file)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", fmt.Errorf("save: %w", err)
	}
	tmp.Close()
	return tmp.Name(), nil
}

func jsonResp(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) apiAddSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jsonResp(w, 405, map[string]string{"error": "POST required"})
		return
	}

	name := r.FormValue("name")
	if name == "" {
		jsonResp(w, 400, map[string]string{"error": "name is required"})
		return
	}

	path, err := s.saveUpload(r)
	if err != nil {
		jsonResp(w, 400, map[string]string{"error": err.Error()})
		return
	}
	defer os.Remove(path)

	err = pipeline.IngestPipeline(s.db, path, name)
	if err != nil {
		jsonResp(w, 500, map[string]string{"error": err.Error()})
		return
	}

	jsonResp(w, 200, map[string]string{"status": "ok", "song": name})
}

func (s *Server) apiMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jsonResp(w, 405, map[string]string{"error": "POST required"})
		return
	}

	path, err := s.saveUpload(r)
	if err != nil {
		jsonResp(w, 400, map[string]string{"error": err.Error()})
		return
	}
	defer os.Remove(path)

	result, err := pipeline.MatchFile(s.db, path)
	if err != nil {
		jsonResp(w, 500, map[string]string{"error": err.Error()})
		return
	}

	if result == "" {
		jsonResp(w, 200, map[string]string{"status": "no match"})
		return
	}

	jsonResp(w, 200, map[string]string{"status": "matched", "match": result})
}

func (s *Server) apiSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := storage.ListSongs(s.db)
	if err != nil {
		jsonResp(w, 500, map[string]string{"error": err.Error()})
		return
	}
	if songs == nil {
		songs = []string{}
	}
	jsonResp(w, 200, map[string]any{"songs": songs})
}

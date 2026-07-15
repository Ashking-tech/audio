package db

import (
	"database/sql"
	"testing"

	"github.com/Ashking-tech/audio/fingerprint"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := InitializeDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestInitializeDB(t *testing.T) {
	db, err := InitializeDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var tableCount int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master
		WHERE type='table' AND name IN ('songs', 'fingerprints')
	`).Scan(&tableCount)
	if err != nil {
		t.Fatal(err)
	}
	if tableCount != 2 {
		t.Errorf("expected 2 tables, got %d", tableCount)
	}
}

func TestInsertsong(t *testing.T) {
	db := setupDB(t)

	fps := []fingerprint.Fingerprint{
		{Hash: 123, AnchorTime: 10},
		{Hash: 456, AnchorTime: 20},
	}

	id, err := Insertsong(db, "test_song", fps)
	if err != nil {
		t.Fatal(err)
	}
	if id == 0 {
		t.Fatal("expected non-zero song ID")
	}

	var name string
	err = db.QueryRow("SELECT name FROM songs WHERE id = ?", id).Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	if name != "test_song" {
		t.Errorf("name = %q, want %q", name, "test_song")
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM fingerprints WHERE song_id = ?", id).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("expected 2 fingerprints, got %d", count)
	}
}

func TestInsertsong_duplicateName(t *testing.T) {
	db := setupDB(t)

	fps := []fingerprint.Fingerprint{{Hash: 1, AnchorTime: 0}}
	_, err := Insertsong(db, "dupe", fps)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Insertsong(db, "dupe", fps)
	if err == nil {
		t.Fatal("expected error for duplicate song name")
	}
}

func TestInsertFingerprints(t *testing.T) {
	db := setupDB(t)

	result, err := db.Exec("INSERT INTO songs (name) VALUES (?)", "test")
	if err != nil {
		t.Fatal(err)
	}
	songID, _ := result.LastInsertId()

	fps := []fingerprint.Fingerprint{
		{Hash: 100, AnchorTime: 5},
		{Hash: 200, AnchorTime: 15},
		{Hash: 300, AnchorTime: 25},
	}

	err = InsertFingerprints(db, songID, fps)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM fingerprints WHERE song_id = ?", songID).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("expected 3 fingerprints, got %d", count)
	}
}

func TestLookUpMatches_exact(t *testing.T) {
	db := setupDB(t)

	fps := []fingerprint.Fingerprint{
		{Hash: 42, AnchorTime: 100},
		{Hash: 99, AnchorTime: 200},
	}

	_, err := Insertsong(db, "song_a", fps)
	if err != nil {
		t.Fatal(err)
	}

	queryFps := []fingerprint.Fingerprint{
		{Hash: 42, AnchorTime: 100},
		{Hash: 99, AnchorTime: 200},
	}

	match, err := LookUpMatches(db, queryFps)
	if err != nil {
		t.Fatal(err)
	}
	if match != "song_a" {
		t.Errorf("match = %q, want %q", match, "song_a")
	}
}

func TestLookUpMatches_noMatch(t *testing.T) {
	db := setupDB(t)

	fps := []fingerprint.Fingerprint{
		{Hash: 42, AnchorTime: 100},
	}

	_, err := Insertsong(db, "song_a", fps)
	if err != nil {
		t.Fatal(err)
	}

	queryFps := []fingerprint.Fingerprint{
		{Hash: 9999, AnchorTime: 0},
	}

	match, err := LookUpMatches(db, queryFps)
	if err != nil {
		t.Fatal(err)
	}
	if match != "" {
		t.Errorf("expected empty match, got %q", match)
	}
}

func TestLookUpMatches_offsetAlignment(t *testing.T) {
	db := setupDB(t)

	fps := []fingerprint.Fingerprint{
		{Hash: 10, AnchorTime: 50},
		{Hash: 20, AnchorTime: 100},
		{Hash: 30, AnchorTime: 150},
	}

	_, err := Insertsong(db, "song_x", fps)
	if err != nil {
		t.Fatal(err)
	}

	queryFps := []fingerprint.Fingerprint{
		{Hash: 10, AnchorTime: 100},
		{Hash: 20, AnchorTime: 150},
		{Hash: 30, AnchorTime: 200},
	}

	match, err := LookUpMatches(db, queryFps)
	if err != nil {
		t.Fatal(err)
	}
	if match != "song_x" {
		t.Errorf("match = %q, want %q", match, "song_x")
	}
}

func TestLookUpMatches_bestMatch(t *testing.T) {
	db := setupDB(t)

	_, err := Insertsong(db, "song_a", []fingerprint.Fingerprint{
		{Hash: 1, AnchorTime: 10},
		{Hash: 2, AnchorTime: 20},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = Insertsong(db, "song_b", []fingerprint.Fingerprint{
		{Hash: 1, AnchorTime: 100},
		{Hash: 2, AnchorTime: 200},
		{Hash: 3, AnchorTime: 300},
	})
	if err != nil {
		t.Fatal(err)
	}

	queryFps := []fingerprint.Fingerprint{
		{Hash: 1, AnchorTime: 10},
		{Hash: 2, AnchorTime: 20},
		{Hash: 3, AnchorTime: 30},
	}

	match, err := LookUpMatches(db, queryFps)
	if err != nil {
		t.Fatal(err)
	}
	if match != "song_a" {
		t.Errorf("expected match 'song_a', got %q", match)
	}
}

func TestListSongs(t *testing.T) {
	db := setupDB(t)

	songs, err := ListSongs(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 0 {
		t.Errorf("expected empty list, got %d songs", len(songs))
	}

	for _, name := range []string{"c", "a", "b"} {
		_, err := Insertsong(db, name, []fingerprint.Fingerprint{
			{Hash: uint32(len(name)), AnchorTime: 0},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	songs, err = ListSongs(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 3 {
		t.Fatalf("expected 3 songs, got %d", len(songs))
	}
	if songs[0] != "c" {
		t.Errorf("songs[0] = %q, want %q", songs[0], "c")
	}
	if songs[1] != "a" {
		t.Errorf("songs[1] = %q, want %q", songs[1], "a")
	}
	if songs[2] != "b" {
		t.Errorf("songs[2] = %q, want %q", songs[2], "b")
	}
}

func TestLookUpMatches_emptyInput(t *testing.T) {
	db := setupDB(t)

	match, err := LookUpMatches(db, []fingerprint.Fingerprint{})
	if err != nil {
		t.Fatal(err)
	}
	if match != "" {
		t.Errorf("expected empty match, got %q", match)
	}
}

func TestLookUpMatches_wrongOffset(t *testing.T) {
	db := setupDB(t)

	_, err := Insertsong(db, "song", []fingerprint.Fingerprint{
		{Hash: 1, AnchorTime: 10},
		{Hash: 2, AnchorTime: 20},
	})
	if err != nil {
		t.Fatal(err)
	}

	queryFps := []fingerprint.Fingerprint{
		{Hash: 1, AnchorTime: 999},
		{Hash: 2, AnchorTime: 888},
	}

	match, err := LookUpMatches(db, queryFps)
	if err != nil {
		t.Fatal(err)
	}
	if match != "song" {
		t.Errorf("expected match 'song', got %q", match)
	}
}

func BenchmarkInsertsong(b *testing.B) {
	db, err := InitializeDB(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	fps := make([]fingerprint.Fingerprint, 1000)
	for i := range fps {
		fps[i] = fingerprint.Fingerprint{
			Hash:       uint32(i),
			AnchorTime: i * 10,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Insertsong(db, "bench", fps)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLookUpMatches(b *testing.B) {
	db, err := InitializeDB(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	fps := make([]fingerprint.Fingerprint, 1000)
	for i := range fps {
		fps[i] = fingerprint.Fingerprint{
			Hash:       uint32(i),
			AnchorTime: i * 10,
		}
	}
	_, err = Insertsong(db, "bench", fps)
	if err != nil {
		b.Fatal(err)
	}

	queryFps := fps[:100]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LookUpMatches(db, queryFps)
		if err != nil {
			b.Fatal(err)
		}
	}
}

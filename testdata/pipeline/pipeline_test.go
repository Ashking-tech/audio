package pipeline_test

import (
	"database/sql"
	"testing"

	"github.com/Ashking-tech/audio/internal/fingerprint"
	"github.com/Ashking-tech/audio/internal/pipeline"
	"github.com/Ashking-tech/audio/internal/storage"
	_ "modernc.org/sqlite"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	database, err := storage.InitializeDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}

func TestMatchRecording_noMatch(t *testing.T) {
	database := setupDB(t)
	samples := fingerprint.GenerateSineWave(440, 0.5, 44100)
	match, err := pipeline.MatchRecording(database, samples)
	if err != nil {
		t.Fatal(err)
	}
	if match != "" {
		t.Errorf("expected empty match, got %q", match)
	}
}

func TestMatchRecording_selfMatch(t *testing.T) {
	database := setupDB(t)

	samples := fingerprint.GenerateSineWave(440, 0.5, 44100)

	spec := fingerprint.Spectogram{WindowSize: 4096, HopSize: 512}
	spectrogram := spec.GenerateSpectogram(samples)
	peaks := fingerprint.FindPeaks(spectrogram, 10, 0.1)
	fps := fingerprint.FingerprintPeaks(peaks, 5)

	_, err := storage.Insertsong(database, "sine_440", fps)
	if err != nil {
		t.Fatal(err)
	}

	match, err := pipeline.MatchRecording(database, samples)
	if err != nil {
		t.Fatal(err)
	}
	if match != "sine_440" {
		t.Errorf("expected match 'sine_440', got %q", match)
	}
}

func TestMatchRecording_onlyCorrectSongMatches(t *testing.T) {
	database := setupDB(t)

	samples := fingerprint.GenerateSineWave(440, 0.5, 44100)

	spec := fingerprint.Spectogram{WindowSize: 4096, HopSize: 512}
	spectrogram := spec.GenerateSpectogram(samples)
	peaks := fingerprint.FindPeaks(spectrogram, 10, 0.1)
	fps := fingerprint.FingerprintPeaks(peaks, 5)

	_, err := storage.Insertsong(database, "target", fps)
	if err != nil {
		t.Fatal(err)
	}

	otherSamples := make([]float64, 44100)
	otherSpec := spec.GenerateSpectogram(otherSamples)
	otherPeaks := fingerprint.FindPeaks(otherSpec, 10, 0.1)
	otherFps := fingerprint.FingerprintPeaks(otherPeaks, 5)
	_, err = storage.Insertsong(database, "silence", otherFps)
	if err != nil {
		t.Fatal(err)
	}

	match, err := pipeline.MatchRecording(database, samples)
	if err != nil {
		t.Fatal(err)
	}
	if match != "target" {
		t.Errorf("expected 'target', got %q", match)
	}
}

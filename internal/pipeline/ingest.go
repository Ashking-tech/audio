package pipeline

import (
	"database/sql"

	"github.com/Ashking-tech/audio/internal/audio"
	"github.com/Ashking-tech/audio/internal/fingerprint"
	"github.com/Ashking-tech/audio/internal/storage"
)

func IngestPipeline(database *sql.DB, path string, songName string) error {
	samples, err := audio.DecodeWav(path)
	if err != nil {
		return err
	}

	spec := fingerprint.Spectogram{WindowSize: 4096, HopSize: 512}
	spectrogram := spec.GenerateSpectogram(samples)

	peaks := fingerprint.FindPeaks(spectrogram, 10, 0.1)

	fps := fingerprint.FingerprintPeaks(peaks, 5)

	_, err = storage.Insertsong(database, songName, fps)
	if err != nil {
		return err
	}
	return nil
}

func MatchFile(database *sql.DB, path string) (string, error) {
	samples, err := audio.DecodeWav(path)

	if err != nil {
		return "", err
	}

	spec := fingerprint.Spectogram{WindowSize: 4096, HopSize: 512}
	spectrogram := spec.GenerateSpectogram(samples)

	peaks := fingerprint.FindPeaks(spectrogram, 10, 0.1)

	fps := fingerprint.FingerprintPeaks(peaks, 5)

	return storage.LookUpMatches(database, fps)
}

func MatchRecording(database *sql.DB, samples []float64) (string, error) {
	spec := fingerprint.Spectogram{WindowSize: 4096, HopSize: 512}
	spectrogram := spec.GenerateSpectogram(samples)

	peaks := fingerprint.FindPeaks(spectrogram, 10, 0.1)

	fps := fingerprint.FingerprintPeaks(peaks, 5)

	return storage.LookUpMatches(database, fps)
}

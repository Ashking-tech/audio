package pipeline

import (
	"database/sql"

	"github.com/Ashking-tech/audio/db"
	"github.com/Ashking-tech/audio/decode"
	"github.com/Ashking-tech/audio/fingerprint"
)


func IngestPipeline(database *sql.DB, path string,songName string) error {
	samples,err := decode.DecodeWav(path)
	if err != nil {
		return err
	}

	spec := fingerprint.Spectogram{WindowSize: 4096,HopSize: 512}
	spectrogram := spec.GenerateSpectogram(samples)

	peaks := fingerprint.FindPeaks(spectrogram, 10, 0.1)

	fps := fingerprint.FingerprintPeaks(peaks, 5)

	_, err = db.Insertsong(database, songName, fps)
	if err != nil {
		return err
	}
	return nil
}


func MatchFile(database *sql.DB,path string)(string,error){
	samples,err := decode.DecodeWav(path)
	
	if err != nil {
		return "",err
	}

	spec := fingerprint.Spectogram{WindowSize: 4096,HopSize: 512}
	spectrogram := spec.GenerateSpectogram(samples)

	peaks := fingerprint.FindPeaks(spectrogram, 10, 0.1)

	fps := fingerprint.FingerprintPeaks(peaks, 5)

	return db.LookUpMatches(database, fps)
}

func MatchRecording(database *sql.DB, samples []float64) (string, error) {
	spec := fingerprint.Spectogram{WindowSize: 4096, HopSize: 512}
	spectrogram := spec.GenerateSpectogram(samples)

	peaks := fingerprint.FindPeaks(spectrogram, 10, 0.1)

	fps := fingerprint.FingerprintPeaks(peaks, 5)

	return db.LookUpMatches(database, fps)
}
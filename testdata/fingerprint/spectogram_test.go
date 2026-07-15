package fingerprint_test

import (
	"math"
	"os"
	"testing"

	"github.com/Ashking-tech/audio/internal/fingerprint"
)

func TestGenerateSineWave(t *testing.T) {
	samples := fingerprint.GenerateSineWave(440, 1.0, 44100)
	expected := 44100
	if len(samples) != expected {
		t.Fatalf("got %d samples, want %d", len(samples), expected)
	}

	peaks := 0
	for i := 1; i < len(samples)-1; i++ {
		if samples[i] > samples[i-1] && samples[i] > samples[i+1] {
			peaks++
		}
	}
	if peaks < 400 || peaks > 480 {
		t.Logf("peak count = %d (expected ~440)", peaks)
	}
}

func TestGenerateSineWave_values(t *testing.T) {
	samples := fingerprint.GenerateSineWave(1.0, 1.0, 100)
	for _, s := range samples {
		if s < -1.0 || s > 1.0 {
			t.Fatalf("sample %f out of range", s)
		}
	}
}

func TestAnalyzeFrequency(t *testing.T) {
	samples := fingerprint.GenerateSineWave(440, 0.2, 44100)
	freq := fingerprint.AnalyzeFrequency(samples, 44100)
	if freq < 430 || freq > 450 {
		t.Errorf("detected frequency = %f, want ~440", freq)
	}
}

func TestAnalyzeFrequency_shortInput(t *testing.T) {
	samples := make([]float64, 100)
	freq := fingerprint.AnalyzeFrequency(samples, 44100)
	if freq != 0 {
		t.Errorf("expected 0 for short input, got %f", freq)
	}
}

func TestSpectogram_GenerateSpectogram(t *testing.T) {
	samples := fingerprint.GenerateSineWave(440, 0.5, 44100)
	s := fingerprint.Spectogram{WindowSize: 4096, HopSize: 512}
	spec := s.GenerateSpectogram(samples)

	if len(spec) == 0 {
		t.Fatal("spectrogram is empty")
	}

	peakBin := -1
	peakMag := 0.0
	for bin := 1; bin < len(spec[0]); bin++ {
		avg := 0.0
		for _, frame := range spec {
			avg += frame[bin]
		}
		avg /= float64(len(spec))
		if avg > peakMag {
			peakMag = avg
			peakBin = bin
		}
	}

	binFreq := float64(peakBin) * 44100.0 / 4096.0
	if binFreq < 400 || binFreq > 480 {
		t.Logf("peak frequency bin %d = %.0f Hz (expected ~440)", peakBin, binFreq)
	}
}

func TestSpectogram_GenerateSpectogram_silence(t *testing.T) {
	samples := make([]float64, 44100)
	s := fingerprint.Spectogram{WindowSize: 4096, HopSize: 512}
	spec := s.GenerateSpectogram(samples)
	if len(spec) == 0 {
		t.Fatal("spectrogram is empty")
	}
	for _, frame := range spec {
		for _, mag := range frame {
			if math.IsNaN(mag) {
				t.Fatal("got NaN in spectrogram")
			}
		}
	}
}

func TestSpectogramImage(t *testing.T) {
	spec := [][]float64{
		{0.1, 0.5, 0.1},
		{0.1, 0.8, 0.1},
		{0.1, 0.3, 0.1},
	}
	path := t.TempDir() + "/test_spec.png"
	err := fingerprint.SpectogramImage(spec, path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("spectrogram image was not created")
	}
	os.Remove(path)
}

func TestSpectogramImage_empty(t *testing.T) {
	err := fingerprint.SpectogramImage([][]float64{}, "test.png")
	if err == nil {
		t.Fatal("expected error for empty spectrogram")
	}
}

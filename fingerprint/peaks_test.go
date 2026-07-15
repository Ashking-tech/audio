package fingerprint

import (
	"testing"
)

func TestFindPeaks_simple(t *testing.T) {
	spec := [][]float64{
		{0.1, 0.2, 0.1},
		{0.1, 0.9, 0.3},
		{0.1, 0.2, 0.1},
	}
	peaks := FindPeaks(spec, 1, 0.3)
	if len(peaks) == 0 {
		t.Fatal("expected at least 1 peak")
	}
	found := false
	for _, p := range peaks {
		if p.TimeBin == 1 && p.FreqBin == 1 {
			found = true
			if p.Magnitude != 0.9 {
				t.Errorf("magnitude = %f, want 0.9", p.Magnitude)
			}
		}
	}
	if !found {
		t.Error("expected peak at (1, 1)")
	}
}

func TestFindPeaks_minMagnitude(t *testing.T) {
	spec := [][]float64{
		{0.5, 0.3},
		{0.3, 0.2},
	}
	peaks := FindPeaks(spec, 1, 0.6)
	if len(peaks) != 0 {
		t.Errorf("expected 0 peaks with high threshold, got %d", len(peaks))
	}
}

func TestFindPeaks_empty(t *testing.T) {
	peaks := FindPeaks([][]float64{}, 5, 0.1)
	if peaks != nil {
		t.Error("expected nil for empty spectrogram")
	}
}

func TestFindPeaks_noDominant(t *testing.T) {
	spec := [][]float64{
		{0.5, 0.5},
		{0.5, 0.5},
	}
	peaks := FindPeaks(spec, 1, 0.1)
	if len(peaks) != 0 {
		t.Errorf("expected 0 peaks for uniform grid, got %d", len(peaks))
	}
}

func TestFindPeaks_edge(t *testing.T) {
	spec := [][]float64{
		{0.8, 0.2},
		{0.2, 0.1},
	}
	peaks := FindPeaks(spec, 1, 0.3)
	if len(peaks) != 1 {
		t.Fatalf("expected 1 peak, got %d", len(peaks))
	}
	if peaks[0].TimeBin != 0 || peaks[0].FreqBin != 0 {
		t.Errorf("expected peak at (0,0), got (%d,%d)", peaks[0].TimeBin, peaks[0].FreqBin)
	}
}

func TestFindPeaks_multiple(t *testing.T) {
	spec := [][]float64{
		{0.9, 0.1, 0.8},
		{0.1, 0.2, 0.1},
		{0.7, 0.1, 0.6},
	}
	peaks := FindPeaks(spec, 1, 0.3)
	if len(peaks) < 3 {
		t.Errorf("expected >=3 peaks, got %d", len(peaks))
	}
}

func BenchmarkFindPeaks(b *testing.B) {
	spec := make([][]float64, 200)
	for t := range spec {
		spec[t] = make([]float64, 150)
		for f := range spec[t] {
			spec[t][f] = 0.1
		}
	}
	spec[100][75] = 0.9
	spec[50][30] = 0.8

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindPeaks(spec, 10, 0.1)
	}
}

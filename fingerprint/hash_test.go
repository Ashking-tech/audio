package fingerprint

import (
	"testing"
)

func makePeaks(times, freqs []int, mags []float64) []Peak {
	peaks := make([]Peak, len(times))
	for i := range peaks {
		peaks[i] = Peak{
			TimeBin:   times[i],
			FreqBin:   freqs[i],
			Magnitude: mags[i],
		}
	}
	return peaks
}

func TestFingerprintPeaks_basic(t *testing.T) {
	peaks := makePeaks(
		[]int{0, 10, 20, 30, 40, 50},
		[]int{5, 10, 15, 20, 25, 30},
		[]float64{0.9, 0.8, 0.7, 0.6, 0.5, 0.4},
	)
	fps := FingerprintPeaks(peaks, 3)
	if len(fps) == 0 {
		t.Fatal("expected at least 1 fingerprint")
	}
}

func TestFingerprintPeaks_fanOut(t *testing.T) {
	peaks := makePeaks(
		[]int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]float64{0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.2, 0.1, 0.0},
	)

	fanOut := 4
	fps := FingerprintPeaks(peaks, fanOut)

	// Anchor 0 pairs with 1,2,3,4 = 4 pairs
	// Anchor 1 pairs with 2,3,4,5 = 4 pairs
	// ...
	// Anchor 9 pairs with none (no next peaks)
	expectedPairs := 0
	for i := 0; i < len(peaks); i++ {
		pairs := fanOut
		if i+fanOut >= len(peaks) {
			pairs = len(peaks) - i - 1
		}
		expectedPairs += pairs
	}

	if len(fps) != expectedPairs {
		t.Errorf("got %d fingerprints, want %d", len(fps), expectedPairs)
	}
}

func TestFingerprintPeaks_deterministic(t *testing.T) {
	peaks := makePeaks(
		[]int{0, 15, 30, 45},
		[]int{3, 7, 11, 15},
		[]float64{0.9, 0.8, 0.7, 0.6},
	)

	fps1 := FingerprintPeaks(peaks, 2)
	fps2 := FingerprintPeaks(peaks, 2)

	if len(fps1) != len(fps2) {
		t.Fatalf("length mismatch: %d vs %d", len(fps1), len(fps2))
	}
	for i := range fps1 {
		if fps1[i].Hash != fps2[i].Hash {
			t.Errorf("hash %d: %d vs %d", i, fps1[i].Hash, fps2[i].Hash)
		}
		if fps1[i].AnchorTime != fps2[i].AnchorTime {
			t.Errorf("anchor time %d: %d vs %d", i, fps1[i].AnchorTime, fps2[i].AnchorTime)
		}
	}
}

func TestFingerprintPeaks_hashFormat(t *testing.T) {
	peaks := makePeaks(
		[]int{0, 10},
		[]int{5, 12},
		[]float64{0.9, 0.8},
	)
	fps := FingerprintPeaks(peaks, 1)
	if len(fps) != 1 {
		t.Fatalf("expected 1 fingerprint, got %d", len(fps))
	}

	fp := fps[0]
	if fp.AnchorTime != 0 {
		t.Errorf("AnchorTime = %d, want 0", fp.AnchorTime)
	}

	// f1 = 5, f2 = 12, dt = 10
	expectedHash := (uint32(5) << 21) | (uint32(12) << 10) | uint32(10)
	if fp.Hash != expectedHash {
		t.Errorf("Hash = %d, want %d", fp.Hash, expectedHash)
	}
}

func TestFingerprintPeaks_empty(t *testing.T) {
	fps := FingerprintPeaks([]Peak{}, 5)
	if len(fps) != 0 {
		t.Errorf("expected 0 fingerprints for empty input, got %d", len(fps))
	}
}

func TestFingerprintPeaks_single(t *testing.T) {
	peaks := makePeaks(
		[]int{0},
		[]int{5},
		[]float64{0.9},
	)
	fps := FingerprintPeaks(peaks, 5)
	if len(fps) != 0 {
		t.Errorf("expected 0 fingerprints for single peak, got %d", len(fps))
	}
}

func TestFingerprintPeaks_zeroFanOut(t *testing.T) {
	peaks := makePeaks(
		[]int{0, 10, 20},
		[]int{5, 10, 15},
		[]float64{0.9, 0.8, 0.7},
	)
	fps := FingerprintPeaks(peaks, 0)
	if len(fps) != 0 {
		t.Errorf("expected 0 fingerprints with fanOut=0, got %d", len(fps))
	}
}

func BenchmarkFingerprintPeaks(b *testing.B) {
	peaks := make([]Peak, 2000)
	for i := range peaks {
		peaks[i] = Peak{
			TimeBin:   i * 5,
			FreqBin:   i % 100,
			Magnitude: float64(i) / 2000,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FingerprintPeaks(peaks, 10)
	}
}

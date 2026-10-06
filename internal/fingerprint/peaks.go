package fingerprint

type Peak struct {
	TimeBin   int
	FreqBin   int
	Magnitude float64
}

// FindPeaks extracts a constellation map from the spectrogram.
// minNeighbourWindow: radius of the neighborhood (in frames/bins) to check for local maxima.
// minMagnitude: fraction (0.0–1.0) of the global max magnitude. Peaks below this are filtered out.
func FindPeaks(spec [][]float64, minNeighbourWindow int, minMagnitude float64) []Peak {
	if len(spec) == 0 {
		return nil
	}

	// ponytail: pre-scan to find max magnitude for relative thresholding.
	// Absolute thresholds fail because FFT magnitudes vary wildly between songs.
	// A relative threshold (fraction of max) adapts to any audio level.
	maxMag := 0.0
	for t := 0; t < len(spec); t++ {
		for f := 0; f < len(spec[t]); f++ {
			if spec[t][f] > maxMag {
				maxMag = spec[t][f]
			}
		}
	}
	threshold := maxMag * minMagnitude

	var peaks []Peak
	for t := 0; t < len(spec); t++ {
		for f := 0; f < len(spec[t]); f++ {
			center := spec[t][f]
			if center < threshold {
				continue
			}

			isPeak := true

			for nt := t - minNeighbourWindow; nt <= t+minNeighbourWindow && isPeak; nt++ {
				if nt < 0 || nt >= len(spec) {
					continue
				}

				for nf := f - minNeighbourWindow; nf <= f+minNeighbourWindow; nf++ {
					if nf < 0 || nf >= len(spec[nt]) {
						continue
					}
					if nt == t && nf == f {
						continue
					}

					if spec[nt][nf] >= center {
						isPeak = false
						break
					}
				}
			}
			if isPeak {
				peaks = append(peaks, Peak{
					TimeBin:   t,
					FreqBin:   f,
					Magnitude: center,
				})
			}
		}
	}
	return peaks
}

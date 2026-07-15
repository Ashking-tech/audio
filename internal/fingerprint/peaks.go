package fingerprint

type Peak struct {
	TimeBin   int
	FreqBin   int
	Magnitude float64
}

func FindPeaks(spec [][]float64, minNeighbourWindow int, minMagnitude float64) []Peak {
	if len(spec) == 0 {
		return nil
	}

	var peaks []Peak
	for t := 0; t < len(spec); t++ {
		for f := 0; f < len(spec[t]); f++ {
			center := spec[t][f]
			if center < minMagnitude {
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

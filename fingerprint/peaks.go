//pseudo code
package fingerprint

type Peak struct {
	TimeBin int  //frame index
	FreqBin int  //freq bin index
	Magnitude float64
}

// Iterate over every cell (t, f) in the spectrogram
// - For each cell, check a rectangular neighborhood of ±timeWindow in time and ±freqWindow in frequency
// - If the center cell has the maximum magnitude in that neighborhood (strictly greater than all neighbors), it's a peak
// - Bounds check: at edges the window shrinks naturally — only compare against cells that exist
// - Return all peaks as a slice

func FindPeaks(spec [][]float64,minNeighbourWindow int) []Peak {

	var peaks []Peak
	for t := 0; t < len(spec); t++ {
		for f := 0; f < len(spec[t]); f++ {
			//spec t f in the cell
			center := spec[t][f]
			isPeak := true

			
			for nt := t - minNeighbourWindow; nt <= t + minNeighbourWindow && isPeak; nt++ {
				//skip invalid time indices
				if nt < 0 || nt >= len(spec){
					continue
				}

				
				for nf := f-minNeighbourWindow; nf <= f+minNeighbourWindow; nf++ {
					//skip invalid frequency indices
					
					if nf < 0 || nf >= len(spec[nt]){
						continue
					}
					//skip the center cell itself
					if nt == t && nf == f{
						continue
					}

					//if any neighbour is greater or equal,
					// the center is not peak
					if spec[nt][nf] >= center {
						isPeak = false
						break
					}

					
				}
			}
			if isPeak {
				peaks = append(peaks,Peak{
					TimeBin:   t,
					FreqBin:   f,
					Magnitude: center,
				})
			}

			
			
		}
	}
	return peaks
}

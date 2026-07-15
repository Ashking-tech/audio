package fingerprint

type Fingerprint struct {
	Hash        uint32
	AnchorTime  int
}

func FingerprintPeaks(peaks []Peak, fanOut int) []Fingerprint {
	var Store []Fingerprint

	for i := 0; i < len(peaks); i++ {

		anchor := peaks[i]

		for j := i + 1; j <= i+fanOut && j < len(peaks); j++ {
			f1, f2, dt := anchor.FreqBin, peaks[j].FreqBin, peaks[j].TimeBin-anchor.TimeBin
			hash := (uint32(f1) << 21) | (uint32(f2) << 10) | uint32(dt)

			Store = append(Store, Fingerprint{
				Hash:       hash,
				AnchorTime: peaks[i].TimeBin,
			})
		}

	}
	return Store

}

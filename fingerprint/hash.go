package fingerprint

type Fingerprint struct {
	Hash        uint32
	AnchorTime  int //time frame of the first peak
}	

func FingerprintPeaks(peaks []Peaks,fanOut int) [] Fingerprint {
	for 
}

// Got it. Here's what to do:
// 1. Create fingerprint/hash.go
// - Add package fingerprint
// - Define a Fingerprint struct with fields: Hash uint32, AnchorTime int
// - Define FingerprintPeaks(peaks []Peak, fanOut int) []Fingerprint
// - Logic: loop over peaks, for each peak pair it with the next fanOut peaks after it. For each pair, compute:
// - f1 = first peak's FreqBin
// - f2 = second peak's FreqBin  
// - dt = second peak's TimeBin minus first peak's TimeBin
// - Pack into uint32 using bit shifts (f1 takes ~11 bits, f2 11 bits, dt 10 bits)
// - Save with AnchorTime = first peak's TimeBin
// 2. Update main.go
// - After FindPeaks, call FingerprintPeaks(peaks, 10)
// - Print the fingerprint count
// That's it. Run go build ./... to check for compilation errors, then go run . to see the fingerprint count.
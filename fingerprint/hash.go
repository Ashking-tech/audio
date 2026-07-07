package fingerprint



type Fingerprint struct {
	Hash        uint32
	AnchorTime  int //time frame of the first peak
}	

func FingerprintPeaks(peaks []Peak,fanOut int) [] Fingerprint {
	//for each peak in the array or list
	// there will be an anchor 
	// we need to look at the next 10 peaks after it,like A is anchor then A->B, A->C, A->D, A->E ... A->K (A pairs with each of the next 10)
	// for each of those peaks:
	// 		record the relationships between anchor and peaks
	// 		save it
	// move to the next peak
	// 
	var Store[] Fingerprint

	for i :=0; i < len(peaks); i++ { //for every peaks in the list
		
		anchor := peaks[i] //this peak is my anchor

		for j := i+1; j <= i+fanOut && j < len(peaks); j++ {
			f1,f2, dt := anchor.FreqBin, peaks[j].FreqBin, peaks[j].TimeBin - anchor.TimeBin
		    hash := (uint32(f1) << 21) | (uint32(f2) << 10) | uint32(dt)

						
						Store = append(Store,Fingerprint{
							Hash: hash,
								AnchorTime: peaks[i].TimeBin,
						})
		}
		
	}
return Store 

} 

// - f1 = first peak's FreqBin
// - f2 = second peak's FreqBin  
// - dt = second peak's TimeBin minus first peak's TimeBin
// - Pack into uint32 using bit shifts (f1 takes ~11 bits, f2 11 bits, dt 10 bits)
// - Save with AnchorTime = first peak's TimeBin
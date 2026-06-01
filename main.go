package main

import (
	"fmt"
	// "github.com/Ashking-tech/audio/decode"
	"github.com/Ashking-tech/audio/fingerprint"
)

func main() {
	// data, err := decode.DecodeWav("output.wav")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("bytes read :", len(data))
	// samples := fingerprint.GenerateSineWave(440,1,48000)	
	// fmt.Println(len(samples))
	// fmt.Println(samples[:10])
	// 
	samples := fingerprint.GenerateSineWave(
		440,   // frequency
		1.0,   // duration (seconds)
		48000, // sample rate
	)

	fmt.Println("Generated samples:", len(samples))

	freq := fingerprint.AnalyzeFrequency(samples, 48000)

	fmt.Println("Returned frequency:", freq)
}



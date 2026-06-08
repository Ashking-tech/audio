package main

import (
	"fmt"
	// "github.com/Ashking-tech/audio/decode"
	"github.com/Ashking-tech/audio/fingerprint"
)

func main() {
	// // data, err := decode.DecodeWav("output.wav")
	// // if err != nil {
	// // 	fmt.Println(err)
	// // 	return
	// // }
	// // fmt.Println("bytes read :", len(data))
	// // samples := fingerprint.GenerateSineWave(440,1,48000)	
	// // fmt.Println(len(samples))
	// // fmt.Println(samples[:10])
	// // 
	// samples := fingerprint.GenerateSineWave(
	// 	440,   // frequency
	// 	1.0,   // duration (seconds)
	// 	48000, // sample rate
	// )

	// fmt.Println("Generated samples:", len(samples))

	// freq := fingerprint.AnalyzeFrequency(samples, 48000)

	// fmt.Println("Returned frequency:", freq)
	// 
	spec := fingerprint.Spectogram{
		WindowSize: 4096,
		HopSize:    512,
	}
	
	signal := fingerprint.GenerateSineWave(
			440,   // frequency
			1.0,   // duration
			44100, // sample rate
		)

	result := spec.GenerateSpectogram(signal)

		fmt.Println("Frames:", len(result))
		fmt.Println("Bins per frame:", len(result[0]))

		err := fingerprint.SpectogramImage(
			result,
			"spectrogram.png",
		)

		if err != nil {
			panic(err)
		}

		fmt.Println("Saved spectrogram.png")

	
	samples := fingerprint.GenerateSineWave(440, 1.0, 48000)
	spectrogram := spec.GenerateSpectogram(samples)
	fmt.Println("Frames:", len(spectrogram))
	fmt.Println("Bins per frame:", len(spectrogram[0]))
	for i := 0; i < 5 && i < len(spectrogram); i++ {
		fmt.Printf("Frame %d max bin: %.0f at index %d\n",
			i, maxInSlice(spectrogram[i]), argmax(spectrogram[i]))
	}
}

func maxInSlice(s []float64) float64 {
	m := 0.0
	for _, v := range s {
		if v > m {
			m = v
		}
	}
	return m
}

func argmax(s []float64) int {
	idx := 0
	m := 0.0
	for i, v := range s {
		if v > m {
			m = v
			idx = i
		}
	}
	return idx
}


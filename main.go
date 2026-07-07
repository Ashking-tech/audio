package main

import (
	"fmt"
	"github.com/Ashking-tech/audio/decode"
	"github.com/Ashking-tech/audio/fingerprint"
)

func main() {
	samples, err := decode.DecodeWav("output.wav")
	if err != nil {
		panic(err)
	}

	spec := fingerprint.Spectogram{
		WindowSize: 4096,
		HopSize:    512,
	}
	result := spec.GenerateSpectogram(samples)

	fmt.Println("Frames:", len(result))
	fmt.Println("Bins per frame:", len(result[0]))

	err = fingerprint.SpectogramImage(result, "spectrogram.png")
	if err != nil {
		panic(err)
	}
	fmt.Println("Saved spectrogram.png")

	peaks := fingerprint.FindPeaks(result, 20)
	fmt.Println("Peaks found:", len(peaks))
}

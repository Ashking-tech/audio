package fingerprint

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/dsp/fourier"
)


func GenerateSineWave(
	frequency float64,
	duration float64,
	sampleRate int,
)[]float64{

	totalSamples := int(float64(sampleRate)*duration)
	
	samples := make([]float64,0,totalSamples)

	for i := 0; i < totalSamples; i++ {
		t:= float64(i)/float64(sampleRate)

		value := math.Sin(2*math.Pi* frequency * t) // example y = sin(2*pie*f*t)

		samples = append(samples,value)
	}

	return samples
}


func AnalyzeFrequency(samples []float64, sampleRate int) float64 {
	const chunkSize = 4096

	if len(samples) < chunkSize {
		return 0
	}

	fft := fourier.NewFFT(chunkSize)
	coeffs := fft.Coefficients(nil, samples[:chunkSize])

	maxMagnitude := 0.0
	maxBin := 0

	// Skip bin 0 (DC component)
	for i := 1; i < len(coeffs)/2; i++ {
		mag := cmplx.Abs(coeffs[i])

		if mag > maxMagnitude {
			maxMagnitude = mag
			maxBin = i
		}
	}

	frequency := float64(maxBin) *
		float64(sampleRate) /
		float64(chunkSize)

	fmt.Println("Strongest bin:", maxBin)
	fmt.Println("Magnitude:", maxMagnitude)
	fmt.Println("Detected frequency:", frequency)

	return frequency
}
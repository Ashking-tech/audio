package fingerprint

import (
	"errors"
	"fmt"
	"image"
	"math"
	"math/cmplx"
	"image/color"
"image/png"
"os"
	"gonum.org/v1/gonum/dsp/fourier"
)


type Spectogram struct {
	WindowSize int 
	HopSize int
}

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

func (s *Spectogram)GenerateSpectogram(signal []float64 ) [][]float64{
	var spectogram [][] float64 //spectogram itself
	windowSize := s.WindowSize // windowSize means the part/chunk of the audio samples , 1 window = 4096 samples
	hopSize := s.HopSize //hopSize means how much the samples overlap

		fft := fourier.NewFFT(windowSize)
	//iterating through the signal
	for start := 0; start+windowSize <= len(signal); start += hopSize {
		//process one frame
		frame := signal[start: start+s.WindowSize] //now frame contains the windowsize
		
		windowed := make([]float64,len(frame))

		for i := range frame{
		n := float64(i)
		N := float64(windowSize - 1)
		coeff := 0.5 * (1 - math.Cos((2*math.Pi *n)/N))
			windowed[i] = frame[i] * coeff
		}

		//run FFT
		coeffs := fft.Coefficients(nil,windowed)


		//compute magnitudes
		magnitudes := make([]float64,len(coeffs))
		for i,c := range coeffs {
		magnitudes[i] = math.Sqrt(real(c)*real(c) + imag(c) * imag(c))
		}

		spectogram = append(spectogram,magnitudes)


		
	}

	return spectogram

	
}

func SpectogramImage(spec[][]float64,filename string) error {
	width := len(spec)

	if width == 0 {
		return errors.New("empty")
		
	}

	height := len(spec[0])

	img:= image.NewGray(
		image.Rect(0,0,width,height))

	var max float64
	for _,frame := range spec {
		for _,mag := range frame {
			if mag > max {
				max = mag
			}
		}
	}


	var maxLog float64

for _, frame := range spec {
    for _, mag := range frame {
        v := math.Log10(1 + mag)

        if v > maxLog {
            maxLog = v
        }
    }
}

for x := 0; x < width; x++ {
    for y := 0; y < height; y++ {

        v := math.Log10(1 + spec[x][y])

        intensity := uint8(
            (v / maxLog) * 255,
        )

        img.SetGray(
            x,
            height-1-y,
            color.Gray{Y: intensity},
        )
    }
}

file, err := os.Create(filename)
if err != nil {
    return err
}
defer file.Close()

return png.Encode(file, img)
}
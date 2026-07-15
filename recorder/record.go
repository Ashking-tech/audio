package recorder

import (
	"github.com/gordonklaus/portaudio"
)

func Record(durationSec int, sampleRate int) ([]float64, error) {
	portaudio.Initialize()
	defer portaudio.Terminate()

	numSamples := durationSec * sampleRate
	buf := make([]float64, numSamples)

	stream, err := portaudio.OpenDefaultStream(1, 0, float64(sampleRate), numSamples, func(in []float64) {
		for i := range buf {
			buf[i] = in[i]
		}
	})
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return nil, err
	}
	if err := stream.Read(); err != nil {
		return nil, err
	}
	return buf, stream.Stop()
}
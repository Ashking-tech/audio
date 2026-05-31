package decode

import (
	"encoding/binary"
	"fmt"
	"os"
)

//struct for the parsed wav data
type Metadata struct {
	Channels      int
	SampleRate    int
	BitsPerSample int
	AudioFormat   int
	DataOffset    int
	DataSize      int
}

func ReadAudioFile(path string) ([]byte, error) {
	audioBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return audioBytes, err
	}

	return audioBytes, err

}

func DecodeWav(path string) ([]float64, error) {

	data, err := ReadAudioFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(data) < 44 {
		return nil, fmt.Errorf("file too small to be valid wav")
	}
	metadata, err := ParseHeader(data)
	if err != nil {
		return nil, err
	}

	pcmSamples, err := ReadPCMSamples(data, metadata)
	if err != nil {
		return nil, err
	}

	monoSamples := StereoToMono(pcmSamples, metadata.Channels)

	normalized := Normalize(monoSamples)

	fmt.Println("PCM Samples:", len(pcmSamples))
	fmt.Println("Mono Samples:", len(monoSamples))
	fmt.Println("First 10:", normalized[:10])

	fmt.Println("bytes read :", len(data))
	fmt.Printf("%+v\n", metadata)
	fmt.Println(normalized[10000:10020])
	return normalized, nil
}

func ParseHeader(data []byte) (Metadata, error) {

	header := Metadata{}

	if len(data) < 44 {
		return header, fmt.Errorf("file too small to be a valid WAV")
	}

	if string(data[0:4]) != "RIFF" {
		return header, fmt.Errorf("not a RIFF file")
	}

	if string(data[8:12]) != "WAVE" {
		return header, fmt.Errorf("not a WAV file")
	}
	header.AudioFormat = int(binary.LittleEndian.Uint16(data[20:22]))
	header.Channels = int(binary.LittleEndian.Uint16(data[22:24])) //channel metatdata is in 22-24 bytes
	header.SampleRate = int(binary.LittleEndian.Uint32(data[24:28])) //sample rate is in 24-28 bytes
	header.BitsPerSample = int(binary.LittleEndian.Uint16(data[34:36])) //bitspersample is in 34-36 bytes

	// Find the data chunk
	pos := 12

	for pos+8 <= len(data) {
		chunkID := string(data[pos : pos+4])
		chunkSize := int(binary.LittleEndian.Uint32(data[pos+4 : pos+8]))

		if chunkID == "data" {
			header.DataOffset = pos + 8
			header.DataSize = chunkSize
			return header, nil
		}

		pos += 8 + chunkSize

		// WAV chunks are word-aligned
		if chunkSize%2 == 1 {
			pos++
		}
	}

	fmt.Printf(
		"Channels: %d\nSample Rate: %d\nBits Per Sample: %d\n",
		header.Channels,
		header.SampleRate,
		header.BitsPerSample,
	)

	return header, fmt.Errorf("data chunk not found")
}


func ReadPCMSamples(data []byte, metadata Metadata) ([]int16, error) {
	if metadata.BitsPerSample != 16 {
		return nil, fmt.Errorf("only 16-bit PCM supported")
	}

	start := metadata.DataOffset
	end := start + metadata.DataSize

	if end > len(data) {
		return nil, fmt.Errorf("invalid data chunk")
	}

	pcmData := data[start:end]

	samples := make([]int16, 0, len(pcmData)/2)

	for i := 0; i+1 < len(pcmData); i += 2 {
		sample := int16(binary.LittleEndian.Uint16(pcmData[i : i+2]))
		samples = append(samples, sample)
	}

	return samples, nil
}


func StereoToMono(samples []int16, channels int) []int16 {
	if channels == 1 {
		return samples
	}

	mono := make([]int16, 0, len(samples)/2)

	for i := 0; i+1 < len(samples); i += 2 {
		left := int32(samples[i])
		right := int32(samples[i+1])

		monoSample := int16((left + right) / 2)
		mono = append(mono, monoSample)
	}

	return mono
}

func Normalize(samples []int16) []float64 {
	normalized := make([]float64, len(samples))

	for i, sample := range samples {
		normalized[i] = float64(sample) / 32768.0
	}

	return normalized
}
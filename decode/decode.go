package decode


import (
	"fmt"
	"os"
	"encoding/binary"
)

//struct for the parsed wav data
type Metadata struct {
	Channels  int
	SampleRate int
	BitsPerSample int
	AudioFormat int
	DataOffset int
	DataSize int
}

func ReadAudioFile(path string) ([]byte,error){
	audioBytes,err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error reading file: %v\n",err)
		return audioBytes,err
	}

	fmt.Printf("succesfully read %d bytes\n",len(audioBytes))
   return audioBytes,err
	

}

func DecodeWav(path string)([]float64,error){
	
	data,err := ReadAudioFile(path)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	

	if len(data) < 44 {
		return nil,fmt.Errorf("file too small to be valid wav")
	}
	metadata, err := ParseHeader(data)
    if err != nil {
		return nil, err
}

	
	fmt.Println("bytes read :",len(data))
	fmt.Printf("%+v\n",metadata)
	return nil,nil
}

func ParseHeader(data []byte)(Metadata,error){

	header := Metadata {}

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
	header.BitsPerSample = int(binary.LittleEndian.Uint16(data[34:36]))//bitspersample is in 34-36 bytes

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
	
	return header,fmt.Errorf("data chunk not found")
	//return 
	// Channels
// Sample Rate
// Bits Per Sample
// Audio Format
// Data Chunk Location
// Data Size
}
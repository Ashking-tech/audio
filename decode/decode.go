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

	metadata := ParseHeader(data)
	
	fmt.Println("bytes read :",len(data))
	fmt.Printf("%+v\n",metadata)
	return nil,nil
}

func ParseHeader(data []byte)Metadata{

	header := Metadata {}

	header.Channels = int(binary.LittleEndian.Uint16(data[22:24])) //channel metatdata is in 22-24 bytes
	header.SampleRate = int(binary.LittleEndian.Uint32(data[24:28])) //sample rate is in 24-28 bytes
	header.BitsPerSample = int(binary.LittleEndian.Uint16(data[34:36]))//bitspersample is in 34-36 bytes

	fmt.Printf(
    "Channels: %d\nSample Rate: %d\nBits Per Sample: %d\n",
    header.Channels,
    header.SampleRate,
    header.BitsPerSample,
)
	
	return header
	//return 
	// Channels
// Sample Rate
// Bits Per Sample
// Audio Format
// Data Chunk Location
// Data Size
}
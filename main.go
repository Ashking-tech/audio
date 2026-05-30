package main

import (
	"fmt"
	"os"
)

func ReadAudioFile(path string) ([]byte,error){
	audioBytes,err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error reading file: %v\n",err)
		return audioBytes,err
	}

	fmt.Print("succesfully read",len(audioBytes),"bytes")
   return audioBytes,err
	

}

func main(){
	data,err := ReadAudioFile("audio.mp3")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("bytes read :",len(data))
}
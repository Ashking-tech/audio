package main

import (
"fmt"	
"github.com/Ashking-tech/audio/decode"
)


func main(){
	data,err := decode.DecodeWav("output.wav")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("bytes read :",len(data))
}
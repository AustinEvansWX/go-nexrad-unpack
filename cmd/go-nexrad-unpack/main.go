package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/roguetechh/go-nexrad-unpack/logger"
	"github.com/roguetechh/go-nexrad-unpack/nexrad"
)

func main() {
	messages, err := nexrad.UnpackMessagesFromChunkFile("../../data/421/20220620-005635-008-I")

	if err != nil {
		logger.Error("Failed to read messages: %v", err)
		return
	}

	//	for _, msg := range messages {
	//		for _, data := range msg.MomentData {
	//			data.MomentData = []float32{}
	//		}
	//	}

	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	err = enc.Encode(messages[0])

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(buffer.String())
}

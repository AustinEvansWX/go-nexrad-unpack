package main

import (
	"github.com/roguetechh/go-nexrad-unpack/logger"
	"github.com/roguetechh/go-nexrad-unpack/nexrad"
)

func main() {
	messages, err := nexrad.UnpackMessagesFromChunkFile("../../data/11/20220622-164844-010-I")

	if err != nil {
		logger.Error("Failed to read messages: %v", err)
		return
	}

	for _, msg := range messages {
		for _, data := range msg.MomentData {
			data.MomentData = []float32{}
		}
		logger.PrettyPrintJSON(msg)
	}
}

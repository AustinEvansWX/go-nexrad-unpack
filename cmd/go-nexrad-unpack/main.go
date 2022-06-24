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

	logger.Info("%d", len(messages))
}

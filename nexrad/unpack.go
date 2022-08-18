package nexrad

import (
	"io/ioutil"

	"github.com/roguetechh/go-nexrad-unpack/bytereader"
	"github.com/roguetechh/go-nexrad-unpack/logger"
)

func UnpackMessagesFromChunkFile(filePath string) ([]*Message, error) {
	chunk, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return UnpackChunk(chunk)
}

func unpackRawMessagesFromChunk(reader *bytereader.Reader) [][]byte {
	messages := [][]byte{}

	for {
		reader.StepForward(12)
		if int(reader.Offset) >= len(reader.Data) {
			break
		}
		messageSize := reader.StaticReadShortUint() * 2
		messages = append(messages, reader.ReadBytes(uint32(messageSize)))
	}

	return messages
}

func UnpackChunk(chunk []byte) ([]*Message, error) {
	data := DecompressBzipChunk(chunk)
	reader := bytereader.NewReader(data)

	messages := []*Message{}
	rawMessages := unpackRawMessagesFromChunk(reader)

	for _, rawMessage := range rawMessages {
		rawMessageReader := bytereader.NewReader(rawMessage)

		message := Message{}

		var err error

		message.MessageHeader, err = ReadMessageHeader(rawMessageReader)

		if err != nil {
			logger.Error("Bad Message Header: %v", err)
			continue
		}

		message.DataHeader, err = ReadDataHeader(rawMessageReader)

		if err != nil {
			logger.Error("Bad Data Header: %v", err)
			continue
		}

		message.VolumeData, err = ReadVolumeData(message.DataHeader, rawMessageReader)

		if err != nil {
			logger.Error("Bad Volume Data: %v", err)
			continue
		}

		message.ElevationData, err = ReadElevationData(message.DataHeader, rawMessageReader)

		if err != nil {
			logger.Error("Bad Elevation Data: %v", err)
			continue
		}

		message.RadialData, err = ReadRadialData(message.DataHeader, rawMessageReader)

		if err != nil {
			logger.Error("Bad Radial Data: %v", err)
			continue
		}

		message.MomentData = map[string]*MomentData{}
		badMomentData := false

		for i, pointer := range message.DataHeader.Pointers {
			if i < 3 || pointer == 0 {
				continue
			}

			momentData, err := ReadMomentData(rawMessageReader, pointer)

			if err != nil {
				logger.Error("Bad Moment Data: %v", err)
				badMomentData = true
				break
			}

			message.MomentData[momentData.DataName] = momentData
		}

		if badMomentData {
			continue
		}

		messages = append(messages, &message)
	}

	return messages, nil
}

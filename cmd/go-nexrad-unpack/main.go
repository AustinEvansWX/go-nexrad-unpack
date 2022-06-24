package main

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"encoding/binary"
	"os"

	"github.com/roguetechh/go-nexrad-unpack/logger"
	"github.com/roguetechh/go-nexrad-unpack/nexrad"
)

type ArchiveHeader struct {
	Format    string
	JulianDay uint32
	Millis    uint32
	Station   string
}

func main() {
	messages, err := nexrad.UnpackMessagesFromChunkFile("data/11/20220622-164844-002-I")

	if err != nil {
		logger.Error("Failed to read messages: %v", err)
		return
	}

	logger.Info("%d", len(messages))
}

func GetArchiveHeader(a *Archive) ArchiveHeader {
	format := a.ReadString(8)
	a.SkipBytes(4)
	day := a.ReadUint()
	milliseconds := a.ReadUint()
	stationId := a.ReadString(4)
	return ArchiveHeader{format, day, milliseconds, stationId}
}

type Archive struct {
	File   *os.File
	Offset uint32
}

func (a *Archive) ReadBytes(size uint32) []byte {
	buffer := make([]byte, size)
	a.File.ReadAt(buffer, int64(a.Offset))
	a.Offset += size
	return buffer
}

func (a *Archive) SkipBytes(size uint32) {
	a.Offset += size
}

func (a *Archive) ReadString(length uint32) string {
	bytes := a.ReadBytes(length)
	return string(bytes)
}

func (a *Archive) ReadUint() uint32 {
	bytes := a.ReadBytes(4)
	return binary.BigEndian.Uint32(bytes)
}

func (a *Archive) ReadBzipBlock() []byte {
	blockSize := a.ReadUint()
	data := a.ReadBytes(blockSize)
	bz := bzip2.NewReader(bytes.NewReader(data))
	scanner := bufio.NewScanner(bufio.NewReader(bz))
	uncompressed := []byte{}
	for scanner.Scan() {
		uncompressed = append(uncompressed, scanner.Bytes()...)
	}
	return uncompressed
}

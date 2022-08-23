package nexrad

import (
	"bytes"
	"compress/bzip2"

	"github.com/roguetechh/go-nexrad-unpack/bytereader"
)

func DecompressBzipChunk(chunk []byte) []byte {
	reader := bytereader.NewReader(chunk)
	blockSize := reader.ReadInt()

	if blockSize < 0 {
		blockSize = -blockSize
	}

	data := reader.ReadBytes(uint32(blockSize))
	uncompressed := make([]byte, blockSize*100)

	bytesRead, _ := bzip2.NewReader(bytes.NewReader(data)).Read(uncompressed)
	uncompressed = uncompressed[:bytesRead]

	return uncompressed
}

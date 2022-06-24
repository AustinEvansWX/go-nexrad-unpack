package nexrad

import (
	"bufio"
	"bytes"
	"compress/bzip2"

	"github.com/roguetechh/go-nexrad-unpack/bytereader"
)

func DecompressBzipChunk(chunk []byte) []byte {
	reader := bytereader.NewReader(chunk)

	blockSize := reader.ReadUint()
	data := reader.ReadBytes(blockSize)

	bz := bzip2.NewReader(bytes.NewReader(data))
	scanner := bufio.NewScanner(bufio.NewReader(bz))

	uncompressed := []byte{}

	for scanner.Scan() {
		uncompressed = append(uncompressed, scanner.Bytes()...)
	}

	return uncompressed
}

package nexrad

import "github.com/roguetechh/go-nexrad-unpack/bytereader"

type MessageHeader struct {
	Size                    uint16
	RDARedundantChannel     uint8
	MessageType             uint8
	IDSequenceNumber        uint16
	JulianDate              uint16
	MillisecondsOfDay       uint32
	NumberOfMessageSegments uint16
	MessageSegmentNumber    uint16
}

func ReadMessageHeader(reader *bytereader.Reader) (*MessageHeader, error) {
	header := MessageHeader{
		reader.ReadShortUint(),
		reader.ReadBytes(1)[0],
		reader.ReadBytes(1)[0],
		reader.ReadShortUint(),
		reader.ReadShortUint(),
		reader.ReadUint(),
		reader.ReadShortUint(),
		reader.ReadShortUint(),
	}

	return &header, header.Validate()
}

func (mh *MessageHeader) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Size", float64(mh.Size), 9, 65535},
		{"RDA Redundant Channel", float64(mh.RDARedundantChannel), 0, 10},
		{"Message Type", float64(mh.MessageType), 1, 31},
		{"ID Sequence Number", float64(mh.IDSequenceNumber), 0, 65535},
		{"Julian Date", float64(mh.JulianDate), 1, 65535},
		{"Milliseconds of Day", float64(mh.MillisecondsOfDay), 0, 86399999},
		{"Number of Message Segments", float64(mh.NumberOfMessageSegments), 1, 65535},
		{"Message Segment Number", float64(mh.MessageSegmentNumber), 1, 65535},
	}

	return validateRanges(rangeChecks)
}

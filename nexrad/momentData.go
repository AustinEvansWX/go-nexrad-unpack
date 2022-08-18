package nexrad

import (
	"github.com/roguetechh/go-nexrad-unpack/bytereader"
)

type MomentData struct {
	DataBlockType                 string
	DataName                      string
	Reserved                      []byte
	DataMomentGateCount           uint16
	DataMomentRange               float32
	DataMomentRangeSampleInterval float32
	Tover                         float32
	SnrThreshold                  float32
	ControlFlags                  uint8
	DataWordSize                  uint8
	Scale                         float32
	Offset                        float32
	MomentData                    []float32
}

func ReadMomentData(reader *bytereader.Reader, pointer uint32) (*MomentData, error) {
	reader.Seek(pointer + MESSAGE_HEADER_SIZE)

	momentData := MomentData{
		DataBlockType:                 reader.ReadString(1),
		DataName:                      reader.ReadString(3),
		Reserved:                      reader.ReadBytes(4),
		DataMomentGateCount:           reader.ReadShortUint(),
		DataMomentRange:               float32(reader.ReadShortUint()) / 1000,
		DataMomentRangeSampleInterval: float32(reader.ReadShortUint()) / 1000,
		Tover:                         float32(reader.ReadShortUint()) / 100,
		SnrThreshold:                  float32(reader.ReadShortUint()) / 1000,
		ControlFlags:                  reader.ReadBytes(1)[0],
		DataWordSize:                  reader.ReadBytes(1)[0],
		Scale:                         reader.ReadFloat(),
		Offset:                        reader.ReadFloat(),
		MomentData:                    []float32{},
	}

	err := momentData.Validate()

	if err != nil {
		return nil, err
	}

	bytesPerWord := uint32(momentData.DataWordSize) / 8

	for i := 0; i < int(momentData.DataMomentGateCount); i++ {
		if bytesPerWord == 1 {
			momentData.MomentData = append(momentData.MomentData, (float32(reader.ReadBytes(1)[0])-momentData.Offset)/momentData.Scale)
		} else {
			momentData.MomentData = append(momentData.MomentData, (float32(reader.ReadShortUint())-momentData.Offset)/momentData.Scale)
		}
	}

	return &momentData, err
}

func (md *MomentData) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Data Moment Gate Count", float64(md.DataMomentGateCount), 0, 1840},
		{"DataMomentRange", float64(md.DataMomentRange), 0, 32.768},
		{"DataMomentRangeSampleInterval", float64(md.DataMomentRangeSampleInterval), 0.25, 4},
		{"Tover", float64(md.Tover), 0, 20},
		{"SnrThreshold", float64(md.SnrThreshold), -12, 20},
		{"Control Flags", float64(md.ControlFlags), 0, 3},
		{"Data Word Size", float64(md.DataWordSize), 8, 16},
		{"Scale", float64(md.Scale), 0, 65535},
		{"Offset", float64(md.Offset), -60.5, 65535},
	}

	return validateRanges(rangeChecks)
}

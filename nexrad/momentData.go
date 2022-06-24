package nexrad

import "github.com/roguetechh/go-nexrad-unpack/bytereader"

type MomentData struct {
	DataBlockType                 string
	DataName                      string
	Reserved                      []byte
	DataMomentGateCount           uint16
	DataMomentRange               uint16
	DataMomentRangeSampleInterval uint16
	Tover                         uint16
	SnrThreshold                  uint16
	ControlFlags                  uint8
	DataWordSize                  uint8
	Scale                         float32
	Offset                        float32
	MomentData                    []float32
	MissingDataPoints             uint32
}

func ReadMomentData(reader *bytereader.Reader) (*MomentData, error) {
	momentData := MomentData{
		reader.ReadString(1),
		reader.ReadString(3),
		reader.ReadBytes(4),
		reader.ReadShortUint(),
		reader.ReadShortUint(),
		reader.ReadShortUint(),
		reader.ReadShortUint(),
		reader.ReadShortUint(),
		reader.ReadBytes(1)[0],
		reader.ReadBytes(1)[0],
		reader.ReadFloat(),
		reader.ReadFloat(),
		[]float32{},
		0,
	}

	err := momentData.Validate()

	if err != nil {
		return nil, err
	}

	bytesPerWord := uint32(momentData.DataWordSize) / 8

	for i := 0; i < int(momentData.DataMomentGateCount); i++ {
		if isNextDataBlock(reader.StaticReadUint()) {
			break
		}

		if bytesPerWord == 1 {
			momentData.MomentData = append(momentData.MomentData, (float32(reader.ReadBytes(1)[0])-momentData.Offset)/momentData.Scale)
		} else {
			reader.StepForward(1)
			if isNextDataBlock(reader.StaticReadUint()) {
				break
			}
			reader.StepBackward(1)
			momentData.MomentData = append(momentData.MomentData, (float32(reader.ReadShortUint())-momentData.Offset)/momentData.Scale)
		}
	}

	momentData.MissingDataPoints += uint32(momentData.DataMomentGateCount) - uint32(len(momentData.MomentData))

	return &momentData, err
}

func isNextDataBlock(id uint32) bool {
	return id == 1146242374 || id == 1146766418 || id == 1146112073 || id == 1146243151 || id == 1145259600
}

func (md *MomentData) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Data Moment Gate Count", float64(md.DataMomentGateCount), 0, 1840},
		{"Control Flags", float64(md.ControlFlags), 0, 3},
		{"Data Word Size", float64(md.DataWordSize), 8, 16},
		{"Scale", float64(md.Scale), 0, 65535},
		{"Offset", float64(md.Offset), -60.5, 65535},
	}

	return validateRanges(rangeChecks)
}

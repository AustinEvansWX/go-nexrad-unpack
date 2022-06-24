package nexrad

import "github.com/roguetechh/go-nexrad-unpack/bytereader"

type ElevationData struct {
	DataBlockType       string
	DataName            string
	Size                uint16
	Atmos               int16
	CalibrationConstant float32
}

func ReadElevationData(dataHeader *DataHeader, reader *bytereader.Reader) (*ElevationData, error) {
	reader.Seek(dataHeader.Pointers[1] + MESSAGE_HEADER_SIZE)

	elevationData := ElevationData{
		reader.ReadString(1),
		reader.ReadString(3),
		reader.ReadShortUint(),
		reader.ReadShortInt(),
		reader.ReadFloat(),
	}

	return &elevationData, elevationData.Validate()
}

func (ed *ElevationData) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Size", float64(ed.Size), 12, 12},
		{"Calibration Constant", float64(ed.CalibrationConstant), -99, 99},
	}

	return validateRanges(rangeChecks)
}

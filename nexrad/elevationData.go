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
	reader.ScanToNonZero()

	elevationData := ElevationData{
		DataBlockType:       reader.ReadString(1),
		DataName:            reader.ReadString(3),
		Size:                reader.ReadShortUint(),
		Atmos:               reader.ReadShortInt(),
		CalibrationConstant: reader.ReadFloat(),
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

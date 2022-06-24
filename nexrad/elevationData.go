package nexrad

import "github.com/roguetechh/go-nexrad-unpack/bytereader"

type ElevationData struct {
	DataBlockType       string
	DataName            string
	Size                uint16
	Atmos               float32
	CalibrationConstant float32
}

func ReadElevationData(dataHeader *DataHeader, reader *bytereader.Reader) (*ElevationData, error) {
	reader.ScanToNonZero()

	elevationData := ElevationData{
		DataBlockType:       reader.ReadString(1),
		DataName:            reader.ReadString(3),
		Size:                reader.ReadShortUint(),
		Atmos:               float32(reader.ReadShortInt()) / 1000,
		CalibrationConstant: reader.ReadFloat(),
	}

	return &elevationData, elevationData.Validate()
}

func (ed *ElevationData) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Size", float64(ed.Size), 12, 12},
		{"Atmos", float64(ed.Atmos), -0.02, -0.002},
		{"Calibration Constant", float64(ed.CalibrationConstant), -99, 99},
	}

	return validateRanges(rangeChecks)
}

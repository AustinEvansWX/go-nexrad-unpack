package nexrad

import (
	"fmt"

	"github.com/roguetechh/go-nexrad-unpack/bytereader"
)

type RadialData struct {
	DataBlockType                 string
	DataName                      string
	Size                          uint16
	UnambiguousRange              uint16
	HorizontalNoiseLevel          float32
	VerticalNoiseLevel            float32
	NyquistVelocity               uint16
	Spares                        []byte
	HorizontalCalibrationConstant float32
	VerticalCalibrationConstant   float32
}

func ReadRadialData(dataHeader *DataHeader, reader *bytereader.Reader) (*RadialData, error) {
	reader.Seek(dataHeader.Pointers[2] + MESSAGE_HEADER_SIZE)

	radialData := RadialData{
		reader.ReadString(1),
		reader.ReadString(3),
		reader.ReadShortUint(),
		reader.ReadShortUint(),
		reader.ReadFloat(),
		reader.ReadFloat(),
		reader.ReadShortUint(),
		reader.ReadBytes(2),
		reader.ReadFloat(),
		reader.ReadFloat(),
	}

	return &radialData, radialData.Validate()
}

func (rd *RadialData) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Size", float64(rd.Size), 28, 28},
		{"Horizontal Noise Level", float64(rd.HorizontalNoiseLevel), -100, -50},
		{"Vertical Noise Level", float64(rd.VerticalNoiseLevel), -100, -50},
		{"Horizontal Calibration Constant", float64(rd.HorizontalCalibrationConstant), -99, 99},
		{"Vertical Calibration Constant", float64(rd.VerticalCalibrationConstant), -99, 99},
	}

	err := validateRanges(rangeChecks)

	if err != nil {
		return err
	}

	if rd.Spares[0] != 0 || rd.Spares[1] != 0 {
		return fmt.Errorf("Invalid Spares | Expected [0, 0] | Got [%d, %d]", rd.Spares[0], rd.Spares[1])
	}

	return nil
}

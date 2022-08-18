package nexrad

import (
	"fmt"

	"github.com/roguetechh/go-nexrad-unpack/bytereader"
)

type RadialData struct {
	DataBlockType                 string
	DataName                      string
	Size                          uint16
	UnambiguousRange              float32
	HorizontalNoiseLevel          float32
	VerticalNoiseLevel            float32
	NyquistVelocity               float32
	Spares                        []byte
	HorizontalCalibrationConstant float32
	VerticalCalibrationConstant   float32
}

func ReadRadialData(dataHeader *DataHeader, reader *bytereader.Reader) (*RadialData, error) {
	reader.Seek(dataHeader.Pointers[2] + MESSAGE_HEADER_SIZE)

	radialData := RadialData{
		DataBlockType:                 reader.ReadString(1),
		DataName:                      reader.ReadString(3),
		Size:                          reader.ReadShortUint(),
		UnambiguousRange:              float32(reader.ReadShortUint()) / 10,
		HorizontalNoiseLevel:          reader.ReadFloat(),
		VerticalNoiseLevel:            reader.ReadFloat(),
		NyquistVelocity:               float32(reader.ReadShortUint()) / 100,
		Spares:                        reader.ReadBytes(2),
		HorizontalCalibrationConstant: reader.ReadFloat(),
		VerticalCalibrationConstant:   reader.ReadFloat(),
	}

	return &radialData, radialData.Validate()
}

func (rd *RadialData) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Size", float64(rd.Size), 28, 28},
		{"Unambiguous Range", float64(rd.UnambiguousRange), 115, 511},
		{"Horizontal Noise Level", float64(rd.HorizontalNoiseLevel), -100, -50},
		{"Vertical Noise Level", float64(rd.VerticalNoiseLevel), -100, -50},
		{"Nyquist Velocity", float64(rd.NyquistVelocity), 8, 35.61},
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

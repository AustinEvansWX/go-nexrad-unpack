package nexrad

import "github.com/roguetechh/go-nexrad-unpack/bytereader"

type VolumeData struct {
	DataBlockType                  string
	DataName                       string
	Size                           uint16
	VersionMajor                   uint8
	VersionMinor                   uint8
	Latitude                       float32
	Longitude                      float32
	SiteHeight                     int16
	FeedhornHeight                 uint16
	CalibrationConstant            float32
	HorizontalShvTxPower           float32
	VerticalShvTxPower             float32
	SystemDifferentialReflectivity float32
	InitialSystemDifferentialPhase float32
	VolumeCoveragePatternNumber    uint16
	ProcessingStatus               uint16
}

func ReadVolumeData(dataHeader *DataHeader, reader *bytereader.Reader) (*VolumeData, error) {
	reader.Seek(dataHeader.Pointers[0] + MESSAGE_HEADER_SIZE)

	volumeData := VolumeData{
		reader.ReadString(1),
		reader.ReadString(3),
		reader.ReadShortUint(),
		reader.ReadBytes(1)[0],
		reader.ReadBytes(1)[0],
		reader.ReadFloat(),
		reader.ReadFloat(),
		reader.ReadShortInt(),
		reader.ReadShortUint(),
		reader.ReadFloat(),
		reader.ReadFloat(),
		reader.ReadFloat(),
		reader.ReadFloat(),
		reader.ReadFloat(),
		reader.ReadShortUint(),
		reader.ReadShortUint(),
	}

	return &volumeData, volumeData.Validate()
}

func (vd *VolumeData) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Size", float64(vd.Size), 52, 52},
		{"Version Major", float64(vd.VersionMajor), 1, 255},
		{"Version Minor", float64(vd.VersionMinor), 0, 255},
		{"Latitude", float64(vd.Latitude), 0, 90},
		{"Longitude", float64(vd.Longitude), -180, 180},
		{"Site Height", float64(vd.SiteHeight), -100, 12000},
		{"Feedhorn Height", float64(vd.FeedhornHeight), 0, 1000},
		{"Calibration Constant", float64(vd.CalibrationConstant), -99, 99},
		{"Horizontal Shv Tx Power", float64(vd.HorizontalShvTxPower), 0, 999.9},
		{"Vertical Shv Tx Power", float64(vd.VerticalShvTxPower), 0, 999.9},
		{"System Differential Reflectivity", float64(vd.SystemDifferentialReflectivity), -7.875, 7.75},
		{"Initial System Differential Phase", float64(vd.InitialSystemDifferentialPhase), 0, 360},
		{"Volume Coverage Pattern Number", float64(vd.VolumeCoveragePatternNumber), 0, 767},
	}

	return validateRanges(rangeChecks)
}

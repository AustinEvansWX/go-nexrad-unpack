package nexrad

import (
	"github.com/roguetechh/go-nexrad-unpack/bytereader"
)

type DataHeader struct {
	RadarIdentifier          string
	CollectionTime           uint32
	JulianDate               uint16
	AzimuthNumber            uint16
	AzimuthAngle             float32
	CompressionIndicator     uint8
	Spare                    uint8
	RadialLength             uint16
	AzimuthResolutionSpacing uint8
	RadialStatus             uint8
	ElevationNumber          uint8
	CutSectorNumber          uint8
	ElevationAngle           float32
	RadialSpotBlankingStatus uint8
	AzimuthIndexingMode      uint8
	DataBlockCount           uint16
	Pointers                 []uint32
}

func ReadDataHeader(reader *bytereader.Reader) (*DataHeader, error) {
	header := DataHeader{
		reader.ReadString(4),   // Radar Identifier
		reader.ReadUint(),      // Collection Time
		reader.ReadShortUint(), // Julian Date
		reader.ReadShortUint(), // Azimuth Number
		reader.ReadFloat(),     // Azimuth Angle
		reader.ReadBytes(1)[0], // CompressionIndicator
		reader.ReadBytes(1)[0], // Spare
		reader.ReadShortUint(), // Radial Length
		reader.ReadBytes(1)[0], // Azimuth Resolution Spacing
		reader.ReadBytes(1)[0], // Radial Status
		reader.ReadBytes(1)[0], // Elevation Number
		reader.ReadBytes(1)[0], // Cut Sector Number
		reader.ReadFloat(),     // Elevation Angle
		reader.ReadBytes(1)[0], // Radial Spot Blanking Status
		reader.ReadBytes(1)[0], // Azimuth Indexing Mode
		reader.ReadShortUint(), // Data Block Count
		[]uint32{
			reader.ReadUint(), // Volume Data Pointer
			reader.ReadUint(), // Elevation Data Pointer
			reader.ReadUint(), // Radial Data Pointer
			reader.ReadUint(), // Moment Data Pointer
			reader.ReadUint(), // Moment Data Pointer
			reader.ReadUint(), // Moment Data Pointer
			reader.ReadUint(), // Moment Data Pointer
			reader.ReadUint(), // Moment Data Pointer
			reader.ReadUint(), // Moment Data Pointer
		},
	}

	return &header, header.Validate()
}

func (dh *DataHeader) Validate() error {
	rangeChecks := []*RangeCheck{
		{"Collection Time", float64(dh.CollectionTime), 0, 86399999},
		{"Julian Date", float64(dh.JulianDate), 1, 65535},
		{"Azimuth Number", float64(dh.AzimuthNumber), 1, 720},
		{"Azimuth Angle", float64(dh.AzimuthNumber), 0, 720},
		{"Compression Indicator", float64(dh.CompressionIndicator), 0, 3},
		{"Radial Length", float64(dh.RadialLength), 9352, 14288},
		{"Azimuth Resolution Spacing", float64(dh.AzimuthResolutionSpacing), 1, 2},
		{"Radial Status", float64(dh.RadialStatus), 0, 132},
		{"Elevation Number", float64(dh.ElevationNumber), 1, 25},
		{"Cut Sector Number", float64(dh.CutSectorNumber), 0, 3},
		{"Elevation Angle", float64(dh.ElevationAngle), -7, 70},
		{"Radial Spot Blanking Status", float64(dh.RadialSpotBlankingStatus), 0, 4},
		{"Azimuth Indexing Mode", float64(dh.AzimuthIndexingMode), 0, 100},
		{"Data Block Count", float64(dh.DataBlockCount), 4, 9},
	}

	return validateRanges(rangeChecks)
}

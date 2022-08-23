package nexrad

import (
	"github.com/roguetechh/go-nexrad-unpack/bytereader"
)

type DataHeader struct {
	RadarIdentifier          string   // ICAO Radar Identifier
	CollectionTime           uint32   // Radial data collection time in milliseconds past midnight GMT
	JulianDate               uint16   // Current Julian date
	AzimuthNumber            uint16   // Radial number within elevation scan
	AzimuthAngle             float32  // Azimuth angle at which radial data was collected
	CompressionIndicator     uint8    // Indicates if message type 31 is compressed and what method of compression is used. The Data Header Block is not compressed.
	Spare                    uint8    // Spare and forces halfword alignment
	RadialLength             uint16   // Uncompressed length of the radial in bytes including the Data Header block length
	AzimuthResolutionSpacing uint8    // Azimuthal spacing between adjacent radials
	RadialStatus             uint8    // Radial Status (e.g. first, last)
	ElevationNumber          uint8    // Elevation number within volume scan
	CutSectorNumber          uint8    // Sector Number within cut
	ElevationAngle           float32  // Elevation angle at which radial radar data was collected
	RadialSpotBlankingStatus uint8    // Spot blanking status for current radial, elevation scan and volume scan
	AzimuthIndexingMode      float32  // Azimuth indexing value (Set if azimuth angle is keyed to constant angles)
	DataBlockCount           uint16   // Number of data blocks
	Pointers                 []uint32 // Data block pointers
}

func ReadDataHeader(reader *bytereader.Reader) (*DataHeader, error) {
	header := DataHeader{
		RadarIdentifier:          reader.ReadString(4),
		CollectionTime:           reader.ReadUint(),
		JulianDate:               reader.ReadShortUint(),
		AzimuthNumber:            reader.ReadShortUint(),
		AzimuthAngle:             reader.ReadFloat(),
		CompressionIndicator:     reader.ReadBytes(1)[0],
		Spare:                    reader.ReadBytes(1)[0],
		RadialLength:             reader.ReadShortUint(),
		AzimuthResolutionSpacing: reader.ReadBytes(1)[0],
		RadialStatus:             reader.ReadBytes(1)[0],
		ElevationNumber:          reader.ReadBytes(1)[0],
		CutSectorNumber:          reader.ReadBytes(1)[0],
		ElevationAngle:           reader.ReadFloat(),
		RadialSpotBlankingStatus: reader.ReadBytes(1)[0],
		AzimuthIndexingMode:      float32(reader.ReadBytes(1)[0]) / 100,
		DataBlockCount:           reader.ReadShortUint(),
		Pointers:                 []uint32{},
	}

	for i := 0; i < int(header.DataBlockCount); i++ {
		header.Pointers = append(header.Pointers, reader.ReadUint())
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
		{"Radial Length", float64(dh.RadialLength), 3824, 14288},
		{"Azimuth Resolution Spacing", float64(dh.AzimuthResolutionSpacing), 1, 2},
		{"Radial Status", float64(dh.RadialStatus), 0, 132},
		{"Elevation Number", float64(dh.ElevationNumber), 1, 25},
		{"Cut Sector Number", float64(dh.CutSectorNumber), 0, 3},
		{"Elevation Angle", float64(dh.ElevationAngle), -7, 70},
		{"Radial Spot Blanking Status", float64(dh.RadialSpotBlankingStatus), 0, 4},
		{"Azimuth Indexing Mode", float64(dh.AzimuthIndexingMode), 0, 1},
		{"Data Block Count", float64(dh.DataBlockCount), 4, 10},
	}

	return validateRanges(rangeChecks)
}

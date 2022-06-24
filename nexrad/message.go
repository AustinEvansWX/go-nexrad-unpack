package nexrad

type Message struct {
	MessageHeader *MessageHeader
	DataHeader    *DataHeader
	VolumeData    *VolumeData
	ElevationData *ElevationData
	RadialData    *RadialData
	MomentData    map[string]*MomentData
	MissingBytes  uint32
}

package psn

// NSEntry represents network statistics.
// Simplied from 'NetDev'.
type NSEntry struct {
	Interface string

	ReceiveBytes    string
	ReceivePackets  string
	TransmitBytes   string
	TransmitPackets string

	// extra fields for sorting
	ReceiveBytesNum    uint64
	ReceivePacketsNum  uint64
	TransmitBytesNum   uint64
	TransmitPacketsNum uint64
}

package psn

// TransportProtocol is tcp, tcp6.
type TransportProtocol int

const (
	TCP TransportProtocol = iota
	TCP6
)

var (
	stringToProtocol = map[string]TransportProtocol{
		"tcp":  TCP,
		"tcp6": TCP6,
	}
	protocolToString = map[TransportProtocol]string{
		TCP:  "tcp",
		TCP6: "tcp6",
	}
)

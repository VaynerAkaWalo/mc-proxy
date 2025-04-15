package packet

type Handshake struct {
	Length   int
	Protocol int
	Hostname string
	Port     int
}

package packet

import (
	"bytes"
	"testing"
)

func TestReadHandshakeRejectsMalformedInput(t *testing.T) {
	tests := map[string][]byte{
		"negative length":        {0xff, 0xff, 0xff, 0xff, 0x0f},
		"oversized length":       {0xff, 0xff, 0x03},
		"truncated packet":       {0x10, 0x00},
		"negative string length": {0x08, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x0f, 0x00},
	}

	for name, input := range tests {
		if _, _, err := ReadHandshake(bytes.NewReader(input)); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestReadHandshakeAcrossFragmentedReads(t *testing.T) {
	input := []byte{0x0b, 0x00, 0xff, 0x05, 0x04, 'a', '.', 'b', 'c', 0x63, 0xdd, 0x01}

	handshake, reply, err := ReadHandshake(&oneByteReader{data: input})
	if err != nil {
		t.Fatal(err)
	}

	if handshake.Hostname != "a.bc" || handshake.Port != 25565 || !bytes.Equal(reply, input) {
		t.Errorf("got %+v %x", handshake, reply)
	}
}

type oneByteReader struct {
	data []byte
}

func (r *oneByteReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, bytes.ErrTooLarge
	}

	n := copy(p[:1], r.data)
	r.data = r.data[n:]
	return n, nil
}

package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const (
	SEGMENT      byte = 0x7F
	CONTINUATION byte = 0x80
	HANDSHAKE    byte = 0x00
)

func ReadVarInt(r io.Reader) (int, []byte, error) {
	var value int32
	var position int
	var buffer []byte

	for {
		var b byte
		err := binary.Read(r, binary.BigEndian, &b)
		if err != nil {
			return 0, nil, err
		}

		value |= int32(b&SEGMENT) << position

		buffer = append(buffer, b)

		if (b & CONTINUATION) == 0 {
			break
		}

		position += 7

		if position >= 35 {
			break
		}
	}

	return int(value), buffer, nil
}

func ReadString(r io.Reader) (string, error) {
	length, _, err := ReadVarInt(r)
	if err != nil {
		log.Println("Error while reading string length", err.Error())
		return "", err
	}

	if length < 0 || length > maxHandshakeLength {
		return "", fmt.Errorf("invalid string length %d", length)
	}

	buffer := make([]byte, length)
	_, err = io.ReadFull(r, buffer)
	if err != nil {
		log.Println("Error while reading string", err.Error())
		return "", err
	}

	return string(buffer[:]), nil
}

func ReadShort(r io.Reader) (uint16, error) {
	var i uint16
	err := binary.Read(r, binary.BigEndian, &i)
	return i, err
}

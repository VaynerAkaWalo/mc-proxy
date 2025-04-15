package packet

import (
	"bytes"
	"errors"
	"io"
	"log"
)

func ReadHandshake(r io.Reader) (Handshake, []byte, error) {
	var bytesToReply []byte

	length, lengthBytes, err := ReadVarInt(r)
	if err != nil {
		log.Println("Error reading packet length", err.Error())
		return Handshake{}, nil, err
	}

	packet := make([]byte, length)
	_, err = r.Read(packet)
	if err != nil {
		log.Println("Error while reading packet", err.Error())
		return Handshake{}, nil, err
	}
	bufferWithPacket := bytes.NewBuffer(packet)

	var packetId int
	var protocol int
	var hostname string
	var port uint16

	packetId, _, err = ReadVarInt(bufferWithPacket)
	if err != nil {
		log.Println("Error occurred while reading packet ID", err.Error())
		return Handshake{}, nil, err
	}

	if byte(packetId)&HANDSHAKE != 0 {
		log.Println("Error this packet is not a handshake")
		return Handshake{}, nil, errors.New("unknown packet type")
	}

	protocol, _, err = ReadVarInt(bufferWithPacket)
	if err != nil {
		log.Println("Error while reading packet protocol", err.Error())
		return Handshake{}, nil, err
	}

	hostname, err = ReadString(bufferWithPacket)
	if err != nil {
		log.Println("Error while reading hostname", err.Error())
		return Handshake{}, nil, err
	}

	port, err = ReadShort(bufferWithPacket)
	if err != nil {
		log.Println("Error while reading port")
		return Handshake{}, nil, err
	}

	bytesToReply = append(bytesToReply, lengthBytes...)
	bytesToReply = append(bytesToReply, packet...)

	handshake := Handshake{
		Length:   length,
		Protocol: protocol,
		Hostname: hostname,
		Port:     int(port),
	}

	return handshake, bytesToReply, nil
}

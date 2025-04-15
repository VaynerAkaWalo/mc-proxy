package routing

import (
	"io"
	"log"
	"mc-proxy/internal/packet"
	"net"
)

type Proxy struct {
	lookup LookupTable
}

func NewProxy(lookupTable LookupTable) Proxy {
	return Proxy{
		lookup: lookupTable,
	}
}

func (p *Proxy) Handle(conn net.Conn) {
	defer conn.Close()

	handshake, bytesToReply, err := packet.ReadHandshake(conn)
	if err != nil {
		log.Println("Error while reading handshake", err.Error())
		return
	}
	log.Printf("Connection with hostname %s and port %d", handshake.Hostname, handshake.Port)

	found, serverAddress := p.lookup.AddressLookup(handshake.Hostname)
	if !found {
		log.Printf("Server for hostname %s not found in lookup table", handshake.Hostname)
		return
	}

	serverConn, erro := p.openServerConnection(serverAddress)
	if erro != nil {
		log.Printf("Error while connecting to server %s", handshake.Hostname)
		return
	}

	ch := make(chan bool)

	serverConn.Write(bytesToReply)

	go proxyPackets(conn, serverConn, ch)
	go proxyPackets(serverConn, conn, ch)

	<-ch
	<-ch
}

func (p *Proxy) openServerConnection(address string) (net.Conn, error) {
	return net.Dial("tcp", address)
}

func proxyPackets(out net.Conn, in net.Conn, c chan bool) {
	defer out.Close()
	defer in.Close()
	io.Copy(in, out)

	c <- true
}

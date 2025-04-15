package proxy

import (
	"fmt"
	"io"
	"log"
	"mc-proxy/internal/packet"
	"mc-proxy/internal/routing"
	"net"
)

type Server struct {
	lookup routing.LookupTable
	addr   string
}

func NewProxyServer(addr string, lookupTable routing.LookupTable) Server {
	return Server{
		addr:   addr,
		lookup: lookupTable,
	}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Println("Failed to start TCP listener")
		return err
	}

	log.Println("Successfully started TCP listener")
	for {
		conn, er := ln.Accept()
		if er != nil {
			fmt.Println("Failed to accept client connection")
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	handshake, bytesToReply, err := packet.ReadHandshake(conn)
	if err != nil {
		log.Println("Error while reading handshake", err.Error())
		return
	}
	log.Printf("Connection with hostname %s and port %d", handshake.Hostname, handshake.Port)

	found, serverAddress := s.lookup.AddressLookup(handshake.Hostname)
	if !found {
		log.Printf("Server for hostname %s not found in lookup table", handshake.Hostname)
		return
	}

	serverConn, erro := s.openServerConnection(serverAddress)
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

func (s *Server) openServerConnection(address string) (net.Conn, error) {
	return net.Dial("tcp", address)
}

func proxyPackets(out net.Conn, in net.Conn, c chan bool) {
	defer out.Close()
	defer in.Close()
	io.Copy(in, out)

	c <- true
}

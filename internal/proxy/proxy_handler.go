package proxy

import (
	"context"
	"io"
	"log/slog"
	"mc-proxy/internal/packet"
	"mc-proxy/internal/routing"
	"net"
)

type Handler struct {
	routingTable *routing.LookupTable
}

func NewProxyHandler(routingTable *routing.LookupTable) *Handler {
	return &Handler{
		routingTable: routingTable,
	}
}

func (h *Handler) Handle(ctx context.Context, cc net.Conn) {
	defer cc.Close()

	handshake, bytesToReply, err := packet.ReadHandshake(cc)
	if err != nil {
		slog.ErrorContext(ctx, "Error while processing handshake packet "+err.Error())
		return
	}

	ctx = context.WithValue(ctx, "server_hostname", handshake.Hostname)

	found, serverAddress := h.routingTable.AddressLookup(handshake.Hostname)
	if !found {
		slog.WarnContext(ctx, "Lookup failed to find matching server")
		return
	}

	sc, err := net.Dial("tcp", serverAddress)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to server "+serverAddress)
		return
	}

	ch := make(chan bool)

	sc.Write(bytesToReply)

	go proxyPackets(cc, sc, ch)
	go proxyPackets(sc, cc, ch)

	<-ch
	<-ch
}

func proxyPackets(out net.Conn, in net.Conn, c chan bool) {
	defer out.Close()
	defer in.Close()

	io.Copy(in, out)

	c <- true
}

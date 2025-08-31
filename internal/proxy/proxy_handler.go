package proxy

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"io"
	"log/slog"
	"mc-proxy/internal/packet"
	"mc-proxy/internal/routing"
	"net"
	"time"
)

const (
	Hostname xctx.ContextKey = "hostname"
	Duration xctx.ContextKey = "duration"
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
	startTime := time.Now()

	handshake, bytesToReply, err := packet.ReadHandshake(cc)
	if err != nil {
		slog.ErrorContext(ctx, "Error while processing handshake packet "+err.Error())
		return
	}

	ctx = context.WithValue(ctx, Hostname, handshake.Hostname)

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

	slog.InfoContext(ctx, "Successfully connected to server")

	ch := make(chan bool)

	sc.Write(bytesToReply)

	go proxyPackets(cc, sc, ch)
	go proxyPackets(sc, cc, ch)

	<-ch
	<-ch

	connectionTime := time.Since(startTime).Milliseconds()
	ctx = context.WithValue(ctx, Duration, connectionTime)

	slog.InfoContext(ctx, "Connection closed")
}

func proxyPackets(out net.Conn, in net.Conn, c chan bool) {
	defer out.Close()
	defer in.Close()

	io.Copy(in, out)

	c <- true
}

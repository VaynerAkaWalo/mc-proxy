package tcp

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"github.com/google/uuid"
	"log/slog"
	"net"
)

const (
	ClientIp xctx.ContextKey = "client_ip"
)

type Server struct {
	addr    string
	handler func(ctx context.Context, conn net.Conn)
}

func NewTCPServer(addr string, handler func(context.Context, net.Conn)) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		slog.Error("Failed to start TCP listener on address " + s.addr)
		return err
	}

	for {
		cc, err := ln.Accept()
		if err != nil {
			slog.Error("Failed to accept client connection")
			continue
		}
		ctx := context.Background()
		ctx = context.WithValue(ctx, ClientIp, cc.RemoteAddr().String())
		ctx = context.WithValue(ctx, xctx.Transaction, uuid.New().String())

		go s.handler(ctx, cc)
	}
}

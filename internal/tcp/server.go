package tcp

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net"
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
		ctx = context.WithValue(ctx, "client_ip", cc.RemoteAddr().String())
		ctx = context.WithValue(ctx, "connection_id", uuid.New().String())

		go s.handler(ctx, cc)
	}
}

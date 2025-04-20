package logging

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	slog.Handler
}

func (ch *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if clientIP, ok := ctx.Value("client_ip").(string); ok {
		r.AddAttrs(slog.String("client_ip", clientIP))
	}

	if connectionId, ok := ctx.Value("connection_id").(string); ok {
		r.AddAttrs(slog.String("connection_id", connectionId))
	}

	if hostname, ok := ctx.Value("server_hostname").(string); ok {
		r.AddAttrs(slog.String("server_hostname", hostname))
	}

	return ch.Handler.Handle(ctx, r)
}

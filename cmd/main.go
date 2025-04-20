package main

import (
	"context"
	"log/slog"
	"mc-proxy/internal/logging"
	"mc-proxy/internal/proxy"
	"mc-proxy/internal/routing"
	"mc-proxy/internal/tcp"
	"net"
	"os"
)

func main() {
	log := slog.New(&logging.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)})
	slog.SetDefault(log)

	log.Info("Application mc-proxy has started")

	lookupTable := routing.NewLookupTable()
	managerClient := routing.ManagerClient{
		Addr: "https://blamedevs.com",
	}

	lookupService := routing.NewLookupService(lookupTable, managerClient)
	lookupService.StartLookupService()

	proxyHandler := proxy.NewProxyHandler(lookupTable)

	server := tcp.NewTCPServer(":25565", proxyHandler.Handle)
	if err := server.ListenAndServe(); err != nil {
		log.Error("Error occurred while listening for client connections")
	}
}

func noOpHandler(ctx context.Context, conn net.Conn) {
	slog.InfoContext(ctx, "Handling client connection")
	conn.Close()
}

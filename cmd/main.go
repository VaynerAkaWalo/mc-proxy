package main

import (
	"github.com/VaynerAkaWalo/go-toolkit/xlog"
	"log/slog"
	"mc-proxy/internal/proxy"
	"mc-proxy/internal/routing"
	"mc-proxy/internal/tcp"
)

func main() {
	slog.SetDefault(slog.New(xlog.NewPreConfiguredHandler(proxy.Hostname, tcp.ClientIp, proxy.Duration)))

	slog.Info("Application mc-proxy has started")

	lookupTable := routing.NewLookupTable()
	managerClient := routing.ManagerClient{
		Addr: "http://mc-server-manager.minecraft-server.svc.cluster.local:8080",
	}

	lookupService := routing.NewLookupService(lookupTable, managerClient)
	lookupService.StartLookupService()

	proxyHandler := proxy.NewProxyHandler(lookupTable)

	server := tcp.NewTCPServer(":25565", proxyHandler.Handle)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error occurred while listening for client connections")
	}
}

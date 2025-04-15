package main

import (
	"fmt"
	"log"
	"mc-proxy/internal/proxy"
	"mc-proxy/internal/routing"
)

func main() {
	fmt.Println("Application mc-proxy has started")

	lookupTable := routing.NewLookupTable()
	managerClient := routing.ManagerClient{
		Addr: "https://blamedevs.com",
	}

	lookupService := routing.NewLookupService(*lookupTable, managerClient)
	lookupService.StartLookupService()

	proxyServer := proxy.NewProxyServer(":25565", *lookupTable)

	err := proxyServer.ListenAndServe()
	if err != nil {
		log.Fatalln("Error while listening for TCP connections", err.Error())
	}
}

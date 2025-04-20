package routing

import (
	"fmt"
	"github.com/VaynerAkaWalo/mc-server-manager/pkg/server"
	"log/slog"
	"time"
)

type LookupService struct {
	lookupTable   *LookupTable
	managerClient ManagerClient
}

func NewLookupService(table *LookupTable, client ManagerClient) *LookupService {
	return &LookupService{
		lookupTable:   table,
		managerClient: client,
	}
}

func (ls *LookupService) StartLookupService() {
	ticker := time.NewTicker(15 * time.Second)

	go func() {
		for {
			ls.updateLookups()
			<-ticker.C
		}
	}()
}

func (ls *LookupService) updateLookups() {
	servers, err := ls.managerClient.ListServers()
	if err != nil {
		slog.Error("Failed to get server from manager")
		return
	}

	newLookupTable := make(map[string]string, len(servers))
	for _, serv := range servers {
		newLookupTable[lookupHostname(serv.Name)] = serverRoute(serv)
	}

	ls.lookupTable.UpdateLookupTable(newLookupTable)

	slog.Info("Lookup update complete, current lookups: " + fmt.Sprint(newLookupTable))
}

func lookupHostname(servAddr string) string {
	return servAddr + ".blamedevs.com"
}

func serverRoute(serv server.Response) string {
	return serv.Name + ":25565"
}

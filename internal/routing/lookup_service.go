package routing

import (
	"github.com/VaynerAkaWalo/mc-server-manager/pkg/server"
	"log"
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

func (sd *LookupService) updateLookups() {
	log.Println("Lookup table update started")
	servers, err := sd.managerClient.ListServers()
	if err != nil {
		log.Println("Failed to get servers from manager")
		return
	}

	newLookupTable := make(map[string]string, len(servers))
	for _, serv := range servers {
		newLookupTable[lookupHostname(serv.Name)] = serverRoute(serv)
	}
	log.Println(newLookupTable)

	sd.lookupTable.UpdateLookupTable(newLookupTable)
}

func lookupHostname(servAddr string) string {
	return servAddr + ".blamedevs.com"
}

func serverRoute(serv server.Response) string {
	return "http://" + serv.Name + ".minecraft-server.svc.cluster.local:25565"
}

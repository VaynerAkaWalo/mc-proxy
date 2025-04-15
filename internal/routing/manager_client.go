package routing

import (
	"encoding/json"
	"github.com/VaynerAkaWalo/mc-server-manager/pkg/server"
	"net/http"
)

type ManagerClient struct {
	Addr string
}

func (c *ManagerClient) ListServers() ([]server.Response, error) {
	resp, err := http.Get(c.Addr + "/mc-server-manager/servers")
	if err != nil {
		return nil, err
	}
	var servers []server.Response

	err = json.NewDecoder(resp.Body).Decode(&servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

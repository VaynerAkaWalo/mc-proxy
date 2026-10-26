package routing

import (
	"encoding/json"
	"fmt"
	"github.com/VaynerAkaWalo/mc-server-manager/pkg/server"
	"net/http"
	"net/url"
)

type ManagerClient struct {
	Addr string
}

func (c *ManagerClient) ListServers() ([]server.Response, error) {
	var servers []server.Response
	cursor := ""

	for {
		page, err := c.listServersPage(cursor)
		if err != nil {
			return nil, err
		}

		servers = append(servers, page.Items...)
		if page.NextCursor == "" {
			return servers, nil
		}
		cursor = page.NextCursor
	}
}

func (c *ManagerClient) listServersPage(cursor string) (server.Collection, error) {
	query := url.Values{"limit": {"100"}}
	if cursor != "" {
		query.Set("cursor", cursor)
	}

	resp, err := http.Get(c.Addr + "/api/v1/servers?" + query.Encode())
	if err != nil {
		return server.Collection{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return server.Collection{}, fmt.Errorf("manager responded with status %d", resp.StatusCode)
	}

	var page server.Collection
	err = json.NewDecoder(resp.Body).Decode(&page)
	return page, err
}

package caddy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hafen/pkg/config"
	"hafen/pkg/db"
	"net/http"
)

// This package provides functions for exposing a local service to a domain via the caddy api.

type Upstreams struct {
	Dial string `json:"dial"`
}

type Handle struct {
	Handler   string      `json:"handler"`
	Upstreams []Upstreams `json:"upstreams"`
}

type Match struct {
	Host []string `json:"host"`
}

type Route struct {
	Handle []Handle `json:"handle"`
	Match  []Match  `json:"match"`
}

type CaddyManager struct {
	queries *db.Queries
	config  *config.Config
}

func NewCaddyManager(config *config.Config, queries *db.Queries) *CaddyManager {
	return &CaddyManager{
		config:  config,
		queries: queries,
	}
}

func (c *CaddyManager) AddProxy(ctx context.Context, proxy db.CreateProxyParams) (*db.Proxy, error) {
	data := Route{
		Handle: []Handle{
			{
				Handler: "reverse_proxy",
				Upstreams: []Upstreams{
					{
						Dial: proxy.Upstream,
					},
				},
			},
		},
		Match: []Match{
			{
				Host: []string{proxy.Match},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("http://%s:%d/config/apps/http/servers/srv0/routes", c.config.Caddy.Host, c.config.Caddy.Port)
	resp, err := http.Post(address, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res, err := c.queries.CreateProxy(ctx, proxy)

	return &res, err
}

func (c *CaddyManager) RemoveProxy(match string) error {
	routes, err := c.GetRoutes()
	if err != nil {
		return err
	}

	for i, route := range routes {
		for _, m := range route.Match {
			if m.Host[0] == match {
				address := fmt.Sprintf("http://%s:%d/config/apps/http/servers/srv0/routes/%d", c.config.Caddy.Host, c.config.Caddy.Port, i)
				req, err := http.NewRequest(http.MethodDelete, address, nil)
				if err != nil {
					return err
				}
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					return err
				}
				if res.StatusCode != 200 {
					return fmt.Errorf("unexpected status code: %d", res.StatusCode)
				}
				err = c.queries.DeleteProxyByMatch(context.Background(), match)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *CaddyManager) GetRoutes() ([]Route, error) {
	address := fmt.Sprintf("http://%s:%d/config/apps/http/servers/srv0/routes", c.config.Caddy.Host, c.config.Caddy.Port)
	resp, err := http.Get(address)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var routes []Route
	err = json.NewDecoder(resp.Body).Decode(&routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}

// adds srv0 to the caddy server
func (c *CaddyManager) LoadDefaultConfig() error {
	proxies, err := c.queries.GetProxies(context.Background())
	if err != nil {
		return err
	}

	routes := make([]Route, 0)
	for _, proxy := range proxies {
		routes = append(routes, Route{
			Handle: []Handle{
				{
					Handler: "reverse_proxy",
					Upstreams: []Upstreams{
						{
							Dial: proxy.Upstream,
						},
					},
				},
			},
			Match: []Match{
				{
					Host: []string{proxy.Match},
				},
			},
		})
	}

	routesData, err := json.Marshal(routes)
	if err != nil {
		return err
	}

	jsonData := []byte(fmt.Sprintf(`
	{
		"apps": {
			"http": {
				"servers": {
					"srv0": {
						"listen": [":443"],
						"routes": %s
					}
				}
			}
		}
	}`, routesData))

	address := fmt.Sprintf("http://%s:%d/config/", c.config.Caddy.Host, c.config.Caddy.Port)
	resp, err := http.Post(address, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

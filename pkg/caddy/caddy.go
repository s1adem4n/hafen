package caddy

import (
	"context"
	"fmt"
	"hafen/pkg/config"
	"hafen/pkg/db"
)

// This package provides functions for exposing a local service to a domain via the caddy api.

type CaddyManager struct {
	queries *db.Queries
	config  *config.Config
	client  *CaddyClient
}

func NewCaddyManager(config *config.Config, queries *db.Queries) *CaddyManager {
	return &CaddyManager{
		config:  config,
		queries: queries,
		client: &CaddyClient{
			ServerName: "srv0",
			Address:    fmt.Sprintf("http://%s:%d", config.Caddy.Host, config.Caddy.Port),
		},
	}
}

func (c *CaddyManager) Init() error {
	err := c.client.Init()
	if err != nil {
		return err
	}

	proxies, err := c.queries.GetProxies(context.Background())
	if err != nil {
		return err
	}

	for _, proxy := range proxies {
		id := GenerateID(proxy.Match)
		if c.client.ObjectExists(fmt.Sprintf("id/%s", id)) {
			continue
		}

		route := NewProxy(
			id,
			proxy.Upstream,
			proxy.Match,
		)
		err := c.client.AddRoute(route)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CaddyManager) AddProxy(ctx context.Context, proxy db.CreateProxyParams) (*db.Proxy, error) {
	route := NewProxy(
		GenerateID(proxy.Match),
		proxy.Upstream,
		proxy.Match,
	)
	err := c.client.AddRoute(route)
	if err != nil {
		return nil, err
	}

	res, err := c.queries.CreateProxy(ctx, proxy)

	return &res, err
}

func (c *CaddyManager) RemoveProxy(proxy db.Proxy) error {
	err := c.client.DeleteObject(fmt.Sprintf("id/%s", GenerateID(proxy.Match)))
	if err != nil {
		return err
	}

	return c.queries.DeleteProxy(context.Background(), proxy.ID)
}

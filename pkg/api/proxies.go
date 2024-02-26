package api

import (
	"hafen/pkg/db"

	"github.com/labstack/echo/v4"
)

func (a *API) GetProxies(c echo.Context) error {
	ctx := c.Request().Context()
	proxies, err := a.queries.GetProxies(ctx)
	if err != nil {
		return err
	}

	return c.JSON(200, proxies)
}

func (a *API) GetProxy(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return RespondError(c, 400, "Invalid ID", err)
	}

	proxy, err := a.queries.GetProxy(ctx, id)
	if err != nil {
		return RespondError(c, 404, "Proxy not found", err)
	}

	return c.JSON(200, proxy)
}

func (a *API) CreateProxy(c echo.Context) error {
	ctx := c.Request().Context()
	var proxy db.CreateProxyParams
	if err := c.Bind(&proxy); err != nil {
		return RespondError(c, 400, "Invalid request", err)
	}

	res, err := a.caddyManager.AddProxy(ctx, proxy)
	if err != nil {
		return RespondError(c, 500, "Failed to create proxy", err)
	}

	return c.JSON(200, res)
}

func (a *API) DeleteProxy(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return RespondError(c, 400, "Invalid ID", err)
	}

	proxy, err := a.queries.GetProxy(ctx, id)
	if err != nil {
		return RespondError(c, 404, "Proxy not found", err)
	}

	err = a.caddyManager.RemoveProxy(proxy.Match)
	if err != nil {
		return RespondError(c, 500, "Failed to delete proxy", err)
	}

	return c.NoContent(200)
}

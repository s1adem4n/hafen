package api

import (
	"hafen/pkg/db"

	"github.com/labstack/echo/v4"
)

func RespondError(c echo.Context, code int, message string, err error) error {
	c.JSON(code, echo.Map{"message": message})
	return err
}

func (a *API) GetTunnels(c echo.Context) error {
	ctx := c.Request().Context()
	tunnels, err := a.queries.GetTunnels(ctx)
	if err != nil {
		return err
	}

	return c.JSON(200, tunnels)
}

func (a *API) GetTunnel(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return RespondError(c, 400, "Invalid ID", err)
	}

	tunnel, err := a.queries.GetTunnel(ctx, id)
	if err != nil {
		return RespondError(c, 404, "Tunnel not found", err)
	}

	return c.JSON(200, tunnel)
}

func (a *API) CreateTunnel(c echo.Context) error {
	ctx := c.Request().Context()
	var tunnel db.CreateTunnelParams
	if err := c.Bind(&tunnel); err != nil {
		return RespondError(c, 400, "Invalid request", err)
	}

	res, err := a.queries.CreateTunnel(ctx, tunnel)
	if err != nil {
		return RespondError(c, 500, "Failed to create tunnel", err)
	}

	return c.JSON(200, res)
}

func (a *API) UpdateTunnel(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return RespondError(c, 400, "Invalid ID", err)
	}
	var tunnel db.UpdateTunnelParams
	if err := c.Bind(&tunnel); err != nil {
		return RespondError(c, 400, "Invalid request", err)
	}

	tunnel.ID = id
	err = a.queries.UpdateTunnel(ctx, tunnel)
	if err != nil {
		return RespondError(c, 500, "Failed to update tunnel", err)
	}

	return c.NoContent(200)
}

func (a *API) DeleteTunnel(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return RespondError(c, 400, "Invalid ID", err)
	}

	err = a.queries.DeleteTunnel(ctx, id)
	if err != nil {
		return RespondError(c, 500, "Failed to delete tunnel", err)
	}

	return c.NoContent(200)
}

func (a *API) StartTunnel(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return RespondError(c, 400, "Invalid ID", err)
	}

	tunnel, err := a.queries.GetTunnel(ctx, id)
	if err != nil {
		return RespondError(c, 404, "Tunnel not found", err)
	}

	if tunnel.Pid != nil {
		return RespondError(c, 400, "Tunnel is already running", nil)
	}

	err = a.tunnelManager.Start(ctx, &tunnel)
	if err != nil {
		return RespondError(c, 500, "Failed to start tunnel", err)
	}

	return c.NoContent(200)
}

func (a *API) StopTunnel(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := ParseInt64(c.Param("id"))
	if err != nil {
		return RespondError(c, 400, "Invalid ID", err)
	}

	tunnel, err := a.queries.GetTunnel(ctx, id)
	if err != nil {
		return RespondError(c, 404, "Tunnel not found", err)
	}

	if tunnel.Pid == nil {
		return RespondError(c, 400, "Tunnel is not running", nil)
	}

	err = a.tunnelManager.Stop(ctx, &tunnel)
	if err != nil {
		return RespondError(c, 500, "Failed to stop tunnel", err)
	}

	return c.NoContent(200)
}

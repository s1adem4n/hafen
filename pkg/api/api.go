package api

import (
	"context"
	"fmt"
	"hafen/frontend"
	"hafen/pkg/caddy"
	"hafen/pkg/config"
	"hafen/pkg/db"
	"hafen/pkg/tunnel"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

type API struct {
	queries       *db.Queries
	config        *config.Config
	echo          *echo.Echo
	tunnelManager *tunnel.TunnelManager
	caddyManager  *caddy.CaddyManager
}

func NewAPI(queries *db.Queries, config *config.Config, tunnelManager *tunnel.TunnelManager, caddyManager *caddy.CaddyManager, logger *slog.Logger) *API {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.CORS())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	return &API{
		queries,
		config,
		e,
		tunnelManager,
		caddyManager,
	}
}

func (a *API) RegisterRoutes() {
	fs := echo.MustSubFS(frontend.Build, "dist")
	a.echo.StaticFS("/", fs)

	a.echo.GET("/tunnels", a.GetTunnels)
	a.echo.GET("/tunnels/:id", a.GetTunnel)
	a.echo.POST("/tunnels", a.CreateTunnel)
	a.echo.PUT("/tunnels/:id", a.UpdateTunnel)
	a.echo.DELETE("/tunnels/:id", a.DeleteTunnel)
	a.echo.POST("/tunnels/:id/start", a.StartTunnel)
	a.echo.POST("/tunnels/:id/stop", a.StopTunnel)

	a.echo.GET("/proxies", a.GetProxies)
	a.echo.GET("/proxies/:id", a.GetProxy)
	a.echo.POST("/proxies", a.CreateProxy)
	a.echo.DELETE("/proxies/:id", a.DeleteProxy)
}

func (a *API) Start() {
	a.echo.Logger.Fatal(a.echo.Start(
		fmt.Sprintf("%s:%d", a.config.API.Host, a.config.API.Port),
	))
}

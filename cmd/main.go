package main

import (
	"context"
	"fmt"
	"hafen/pkg/api"
	"hafen/pkg/caddy"
	"hafen/pkg/config"
	"hafen/pkg/db"
	"hafen/pkg/tunnel"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	config, err := config.ParseConfig()
	if err != nil {
		logger.Error("Failed to parse config", "err", err)
		os.Exit(1)
	}

	if config.Caddy.TunnelServer {
		cmd := exec.Command(
			"ssh",
			"-N", "-T",
			"-o", "ServerAliveInterval=60",
			"-o", "ServerAliveCountMax=10",
			"-o", "ExitOnForwardFailure=yes",
			"-p", fmt.Sprintf("%d", config.Server.Port),
			"-L", "2019:localhost:2019",
			fmt.Sprintf("%s@%s", config.Server.User, config.Server.Host),
		)
		err := cmd.Start()
		if err != nil {
			logger.Error("Failed to start caddy tunnel", "err", err)
			os.Exit(1)
		}
		time.Sleep(1 * time.Second)

		go func() {
			err := cmd.Wait()
			if err != nil {
				logger.Error("Caddy tunnel exited", "err", err)
			}
		}()
	}

	ctx := context.Background()
	conn, err := db.NewConnection(ctx)
	if err != nil {
		logger.Error("Failed to connect to database", "err", err)
		os.Exit(1)
	}

	queries := db.New(conn)
	tunnelManager := tunnel.NewTunnelManager(config, queries)
	caddyManager := caddy.NewCaddyManager(config, queries)

	err = caddyManager.LoadDefaultConfig()
	if err != nil {
		logger.Error("Failed to load default caddy config", "err", err)
		os.Exit(1)
	}

	api := api.NewAPI(queries, config, tunnelManager, caddyManager, logger)
	api.RegisterRoutes()

	slog.Info("Starting server", "host", config.API.Host, "port", config.API.Port)
	go api.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	logger.Info("Stopping tunnels")
	tunnels, err := queries.GetTunnels(ctx)
	if err != nil {
		logger.Error("Failed to get tunnels", "err", err)
	}

	for _, tunnel := range tunnels {
		if tunnel.Pid != nil {
			err = queries.UpdateTunnelPid(ctx, db.UpdateTunnelPidParams{
				Pid: nil,
				ID:  tunnel.ID,
			})
			if err != nil {
				logger.Error("Failed to update tunnel pid", "err", err)
			}
		}
	}

	logger.Info("Shutting down")
}

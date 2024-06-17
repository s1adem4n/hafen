package tunnel

import (
	"context"
	"errors"
	"fmt"
	"hafen/pkg/config"
	"hafen/pkg/db"
	"log/slog"
	"os/exec"
	"strings"
)

type TunnelManager struct {
	config  *config.Config
	queries *db.Queries
}

func NewTunnelManager(config *config.Config, queries *db.Queries) *TunnelManager {
	return &TunnelManager{
		config:  config,
		queries: queries,
	}
}

func (t *TunnelManager) Start(ctx context.Context, tunnel *db.Tunnel) error {
	auth := fmt.Sprintf("%s@%s", t.config.Server.User, t.config.Server.Host)
	// 0.0.0.0 to bind to all interfaces on the server
	tunnelArg := fmt.Sprintf("0.0.0.0:%d:%s:%d", tunnel.RemotePort, tunnel.LocalHost, tunnel.LocalPort)

	args := []string{
		"-N", "-T",
		"-o", "ServerAliveInterval=60",
		"-o", "ServerAliveCountMax=10",
		"-o", "ExitOnForwardFailure=yes",
		"-p", fmt.Sprintf("%d", t.config.Server.Port),
		"-R", tunnelArg,
		auth,
	}
	slog.Info("Starting tunnel", "args", strings.Join(args, " "))
	cmd := exec.Command("ssh", args...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	pid := int64(cmd.Process.Pid)
	err = t.queries.UpdateTunnelPid(ctx, db.UpdateTunnelPidParams{
		Pid: &pid,
		ID:  tunnel.ID,
	})

	go func() {
		cmd.Wait()
		t.queries.UpdateTunnelPid(ctx, db.UpdateTunnelPidParams{
			Pid: nil,
			ID:  tunnel.ID,
		})
	}()

	return err
}

func (t *TunnelManager) Stop(ctx context.Context, tunnel *db.Tunnel) error {
	if tunnel.Pid == nil {
		return errors.New("tunnel is not running")
	}

	// gracefully close the tunnel by sending a SIGTERM
	cmd := exec.Command("kill", "-TERM", fmt.Sprintf("%d", *tunnel.Pid))
	err := cmd.Run()
	if err != nil {
		return err
	}

	err = t.queries.UpdateTunnelPid(ctx, db.UpdateTunnelPidParams{
		Pid: nil,
		ID:  tunnel.ID,
	})

	return err
}

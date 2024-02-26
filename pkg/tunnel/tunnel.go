package tunnel

import (
	"context"
	"errors"
	"fmt"
	"hafen/pkg/config"
	"hafen/pkg/db"
	"os/exec"
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
	tunnelArg := fmt.Sprintf("%d:%s:%d", tunnel.RemotePort, tunnel.LocalHost, tunnel.LocalPort)

	args := []string{
		"-N", "-R",
		tunnelArg,
		"-p", fmt.Sprintf("%d", t.config.Server.Port),
		auth,
	}
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
	cmd := exec.Command("kill", "-INT", fmt.Sprintf("%d", *tunnel.Pid))
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

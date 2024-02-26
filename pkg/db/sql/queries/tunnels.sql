-- name: GetTunnels :many
SELECT *
FROM tunnels;

-- name: GetTunnel :one
SELECT *
FROM tunnels
WHERE id = ?;

-- name: CreateTunnel :one
INSERT INTO tunnels (
    name,
    remote_port,
    local_host,
    local_port,
    pid
  )
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateTunnel :exec
UPDATE tunnels
SET name = ?,
  remote_port = ?,
  local_host = ?,
  local_port = ?,
  pid = ?
WHERE id = ?;

-- name: UpdateTunnelPid :exec
UPDATE tunnels
SET pid = ?
WHERE id = ?;

-- name: DeleteTunnel :exec
DELETE FROM tunnels
WHERE id = ?;
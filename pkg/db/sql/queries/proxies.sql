-- name: GetProxies :many
SELECT *
FROM proxies;

-- name: GetProxy :one
SELECT *
FROM proxies
WHERE id = ?;

-- name: GetProxyByMatch :one
SELECT *
FROM proxies
WHERE match = ?;

-- name: CreateProxy :one
INSERT INTO proxies (name, upstream, match)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateProxy :exec
UPDATE proxies
SET name = ?,
  upstream = ?,
  match = ?
WHERE id = ?;

-- name: DeleteProxy :exec
DELETE FROM proxies
WHERE id = ?;

-- name: DeleteProxyByMatch :exec
DELETE FROM proxies
WHERE match = ?;
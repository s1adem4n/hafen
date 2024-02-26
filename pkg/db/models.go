// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import ()

type Proxy struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Upstream string `json:"upstream"`
	Match    string `json:"match"`
}

type Tunnel struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	RemotePort int64  `json:"remotePort"`
	LocalHost  string `json:"localHost"`
	LocalPort  int64  `json:"localPort"`
	Pid        *int64 `json:"pid"`
}

package frontend

import (
	"embed"
)

//go:embed dist/*
var Build embed.FS

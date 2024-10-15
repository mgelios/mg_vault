package main

import (
	"embed"
	"log/slog"
	"mg_vault/router"
)

//go:embed templates
var templatesFolder embed.FS

//go:embed all:static
var staticContentFolder embed.FS

func main() {
	slog.Info("starting web server")
	router.InitServer(templatesFolder, staticContentFolder)
	router.RunServer()
}

package main

import (
	"embed"
	"mg_vault/router"
)

//go:embed templates
var templatesFolder embed.FS

//go:embed static
var staticContentFolder embed.FS

func main() {
	router.RunServer(templatesFolder)
}

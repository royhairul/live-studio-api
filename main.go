package main

import (
	"github.com/royhairul/live-studio-api/internal/server"
)

func main() {
	app := server.NewApp()
	app.Run()
}

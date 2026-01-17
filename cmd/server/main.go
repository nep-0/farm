package main

import (
	"farm/internal/server"
	"log/slog"
	"os"
)

func main() {
	app, err := server.New("config.json")
	if err != nil {
		slog.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}

	if err := app.Start(); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

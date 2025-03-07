package main

import (
	"posts/internal/app"
	"posts/internal/pkg/config"
)

func main() {
	cfg := config.Load()

	app.Run(cfg)
}

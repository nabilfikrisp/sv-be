// Package main is the application entry point.
package main

import (
	"log"

	"github.com/nabilfikrisp/sv-be/config"
	"github.com/nabilfikrisp/sv-be/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log.Printf("Running on port %s", cfg.HTTP.Port)

	// Run
	app.Run(cfg)
}

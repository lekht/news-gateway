package main

import (
	"log"

	"github.com/lekht/news-gateway/config"
	"github.com/lekht/news-gateway/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}

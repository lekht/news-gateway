package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lekht/news-gateway/config"
	"github.com/lekht/news-gateway/internal/api"
	"github.com/lekht/news-gateway/pkg/server"
)

func New() {}

func Run(cfg *config.Config) {
	api := api.New(&cfg.API)
	router := api.Router()
	httpServer := server.New(router, server.Port(cfg.Server.Listen))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println(fmt.Errorf("app - Run - signal: " + s.String()))
	case err := <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - server.Notify: %w", err))
	}

	err := httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - server.Shutdown: %w", err))
	}
}

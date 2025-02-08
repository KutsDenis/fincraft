package main

import (
	"fincraft/internal/transport/handlers"
	"log"

	"github.com/KutsDenis/logzap"

	"fincraft/internal/config"
	"fincraft/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logzap.Init(cfg.AppEnv)
	defer logzap.Sync()

	handler := handlers.NewHandler()

	s := server.NewServer(cfg.HTTPPort, handler.RegisterRoutes())
	s.Start()
	s.GracefulShutdown()
}

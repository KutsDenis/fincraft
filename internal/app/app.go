package app

import (
	"log"

	"github.com/KutsDenis/logzap"
	"go.uber.org/zap"

	"fincraft/internal/config"
	"fincraft/internal/infra/db"
	"fincraft/internal/transport/handlers"
	"fincraft/internal/transport/server"
)

// Run инициализирует и запускает приложение
func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logzap.Init(cfg.AppEnv)
	defer logzap.Sync()

	initDB(cfg)

	handler := handlers.NewHandler()

	s := server.NewServer(cfg.HTTPPort, handler.RegisterRoutes())

	s.Start()
	s.GracefulShutdown()
}

func initDB(cfg config.Config) {
	database, err := db.Connect(cfg.DBPath)
	if err != nil {
		logzap.Fatal("failed to connect to database", zap.Error(err))
	}
	// noinspection GoUnhandledErrorResult
	defer database.Close()

	if err := db.ApplyMigrations(database, config.MigrationsPath, cfg.AppEnv); err != nil {
		logzap.Fatal("failed to apply migrations:", zap.Error(err))
	}
}

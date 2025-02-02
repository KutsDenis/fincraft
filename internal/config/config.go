package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const envPath = "config/.env"

// Config конфигурация приложения
type Config struct {
	FirstParam  string `env:"FIRST_PARAM"`
	SecondParam string `env:"SECOND_PARAM"`
}

// Load загружает конфигурацию приложения
func Load() (Config, error) {
	err := setEnvFromFile()
	if err != nil {
		return Config{}, err
	}

	cfg, err := loadFromEnv()
	if err != nil {
		return Config{}, nil
	}

	err = validate(cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// setEnvFromFile устанавливает переменные окружения из конфигурационного файла при его наличии
func setEnvFromFile() error {
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	return nil
}

// loadFromEnv загружает конфигурацию из переменных окружения
func loadFromEnv() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// validate валидирует конфигурацию
func validate(cfg Config) error {
	if cfg.FirstParam == "" {
		return fmt.Errorf("first param is required")
	}

	return nil
}

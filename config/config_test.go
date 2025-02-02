package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	filePath := "../.env" // Тестовый файл
	requiredParams := []string{"TEST_PARAM"}

	os.Setenv("TEST_PARAM", "123") // Устанавливаем тестовую переменную окружения

	err := LoadConfig(filePath, requiredParams)
	if err != nil {
		t.Fatalf("LoadConfig() вернул ошибку: %v", err)
	}

	if value, _ := os.LookupEnv("TEST_PARAM"); value != "123" {
		t.Fatalf("Переменная TEST_PARAM не загружена корректно")
	}
}

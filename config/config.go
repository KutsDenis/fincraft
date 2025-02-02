package config

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

func LoadConfig(filePath string, requiredParams []string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return errors.New("ошибка в строке: " + line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			return errors.New("пустое значение у переменной " + key)
		}
		os.Setenv(key, value)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	err = checkRequiredParams(requiredParams)
	if err != nil {
		return err
	}

	log.Println("Конфигурация успешно загружена")
	return nil
}

func checkRequiredParams(requiredParams []string) error {
	var missingVars []string
	for _, param := range requiredParams {
		if strings.TrimSpace(os.Getenv(param)) == "" {
			missingVars = append(missingVars, param)
		}
	}
	if len(missingVars) > 0 {
		return errors.New("отсутствуют обязательные переменные окружения: " + strings.Join(missingVars, ", "))
	}
	return nil
}

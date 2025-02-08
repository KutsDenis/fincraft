.PHONY: test mocks init

# Команда для генерации моков перед запуском тестов
mocks:
	go generate ./...

# Команда для запуска тестов с предварительной генерацией моков
test: mocks
	go test ./... -v

# Команда инициализации проекта
init:
	bash scripts/init.sh && pre-commit install

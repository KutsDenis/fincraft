package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/KutsDenis/logzap"
	"go.uber.org/zap"

	"fincraft/internal/config"
)

const fileSuffix = ".sql"

// ApplyMigrations применяет миграции в зависимости от окружения.
func ApplyMigrations(db *sql.DB, migrationsPath, env string) error {
	logzap.Info("Checking for database migrations...")

	if env == config.ProdEnv {
		return applyProdMigrations(db, migrationsPath)
	}

	return applyDevMigrations(db, migrationsPath)
}

// applyProdMigrations применяет миграции из единого файла (используется в продакшене).
//
// В данном кейсе нет таблицы для хранения информации о примененных миграциях,
// потому что обновления приходят с сервера, и приложение получает только
// нужный пакет миграций, которые необходимы от текущей до версии обновления.
func applyProdMigrations(db *sql.DB, migrationFile string) error {
	logzap.Info("Applying production migrations", zap.String("file", migrationFile))

	migrationSQL, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		return fmt.Errorf("failed to apply production migrations: %w", err)
	}

	logzap.Info("Production migrations applied successfully")

	return nil
}

// applyDevMigrations применяет миграции из директории (используется в локальной разработке).
func applyDevMigrations(db *sql.DB, migrationsDir string) error {
	logzap.Info("Applying development migrations", zap.String("dir", migrationsDir))

	migrationFiles, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return err
	}

	if err := createMigrationTable(db); err != nil {
		return err
	}

	return applyMigrations(db, migrationsDir, migrationFiles)
}

// getMigrationFiles получает список файлов миграций и сортирует их.
func getMigrationFiles(migrationsDir string) ([]string, error) {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string

	for _, file := range files {
		if strings.HasSuffix(file.Name(), fileSuffix) {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	return migrationFiles, nil
}

// createMigrationTable создаёт таблицу для отслеживания миграций, если её нет.
func createMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY
	)`)

	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	return nil
}

// applyMigrations применяет список миграций
func applyMigrations(db *sql.DB, migrationsDir string, migrationFiles []string) error {
	for _, file := range migrationFiles {
		version := strings.TrimSuffix(file, fileSuffix)

		applied, err := isMigrationApplied(db, version)
		if err != nil {
			return err
		}

		if !applied {
			logzap.Info("Applying migration", zap.String("version", version))

			if err := executeMigration(db, migrationsDir, file, version); err != nil {
				return err
			}
		}
	}

	logzap.Info("Development migrations applied successfully")

	return nil
}

// isMigrationApplied проверяет, была ли применена миграция.
func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)", version).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("failed to check migration status for %s: %w", version, err)
	}

	return exists, nil
}

// executeMigration выполняет миграцию и записывает её в БД.
func executeMigration(db *sql.DB, migrationsDir, file, version string) error {
	migrationSQL, err := os.ReadFile(filepath.Join(migrationsDir, file))
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", file, err)
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		return fmt.Errorf("failed to apply migration %s: %w", version, err)
	}

	if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
		return fmt.Errorf("failed to record migration %s: %w", version, err)
	}

	return nil
}

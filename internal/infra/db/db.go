package db

import (
	"database/sql"
	"path/filepath"

	"github.com/KutsDenis/logzap"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const driver = "sqlite3"
const nameDB = "fincraft.db"

// Connect открывает соединение с базой данных.
func Connect(dbPath string) (*sql.DB, error) {
	fullPath := getFullPath(dbPath)

	logzap.Info("connecting to database", zap.String("dbPath", fullPath))

	db, err := sql.Open(driver, fullPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getFullPath возвращает полный путь к базе данных
func getFullPath(dbPath string) string {
	if dbPath == "" {
		dbPath = "./"
	}

	return filepath.Join(dbPath, nameDB)
}

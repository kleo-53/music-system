package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	psql "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	"github.com/kleo-53/music-system/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDBIfNotExists(ctx context.Context, adminDBURL, dbName string) error {
	logger.Log().Debug(ctx, "Connecting to admin database with URL: %s", adminDBURL)
	conn, err := gorm.Open(postgres.Open(adminDBURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to admin database: %v", err)
	}
	db, err := conn.DB()
	if err != nil {
		logger.Log().Info(ctx, "Error getting underlying database connection: %s", err.Error())
		return err
	}
	defer db.Close()

	// Проверяем, существует ли база данных
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	if err := conn.Raw(query).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check if database exists: %v", err)
	}

	if !exists {
		// Создаем базу данных
		createQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
		if err := conn.Exec(createQuery).Error; err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		logger.Log().Info(ctx, "Database %s created successfully.", dbName)
	} else {
		logger.Log().Info(ctx, "Database %s already exists.", dbName)
	}
	return nil
}

func RunMigration(dataBaseURL string) error {
	if err := godotenv.Load("config.env"); err != nil {
		logger.Log().Warn(context.Background(), "No .env file found, using environment variables")
	}
	db, err := sql.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_URL"))
	if err != nil {
		return err
	}

	driver, err := psql.WithInstance(db, &psql.Config{})
	if err != nil {
		return err
	}
	migrationPath := os.Getenv("MIGRATION_PATH")
	if !strings.HasPrefix(migrationPath, "file://") {
		absPath, err := filepath.Abs(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		migrationPath = "file://" + filepath.ToSlash(absPath)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		os.Getenv("DB_TYPE"),
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}

	return nil
}

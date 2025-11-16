package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pavanrkadave/homies/pkg/logger"
	"go.uber.org/zap"
)

// RunMigrations runs database migrations using golang-migrate
func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("Failed to create migration driver", zap.Error(err))
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		logger.Error("Failed to create migrate instance", zap.Error(err))
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logger.Warn("No migration version found (first run)", zap.Error(err))
		version = 0
	}

	if dirty {
		logger.Warn("Database is in dirty state", zap.Uint("version", version))
		return fmt.Errorf("database is in dirty state at version %d", version)
	}

	logger.Info("Current migration version", zap.Uint("version", version))

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("No new migrations to apply")
			return nil
		}
		logger.Error("Failed to run migrations", zap.Error(err))
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Get new version
	newVersion, _, _ := m.Version()
	logger.Info("Migrations applied successfully",
		zap.Uint("from_version", version),
		zap.Uint("to_version", newVersion))

	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	version, _, _ := m.Version()
	logger.Info("Rolling back migration", zap.Uint("current_version", version))

	if err := m.Steps(-1); err != nil {
		logger.Error("Failed to rollback migration", zap.Error(err))
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	newVersion, _, _ := m.Version()
	logger.Info("Migration rolled back successfully",
		zap.Uint("from_version", version),
		zap.Uint("to_version", newVersion))

	return nil
}

// MigrateToVersion migrates to a specific version
func MigrateToVersion(db *sql.DB, migrationsPath string, targetVersion uint) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	currentVersion, _, _ := m.Version()
	logger.Info("Migrating to specific version",
		zap.Uint("current", currentVersion),
		zap.Uint("target", targetVersion))

	if err := m.Migrate(targetVersion); err != nil {
		logger.Error("Failed to migrate to version", zap.Uint("version", targetVersion), zap.Error(err))
		return fmt.Errorf("failed to migrate to version %d: %w", targetVersion, err)
	}

	logger.Info("Successfully migrated to version", zap.Uint("version", targetVersion))
	return nil
}

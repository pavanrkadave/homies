package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	log.Println("Running database migrations...")

	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	sort.Strings(files)

	for _, file := range files {
		log.Printf("Executing migration: %s", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration file %s: %w", file, err)
		}

		log.Printf("✓ Migration completed: %s", filepath.Base(file))
	}
	log.Println("Done running database migrations.")
	return nil
}

package main

import (
	"database/sql"
	"log"

	"github.com/pavanrkadave/homies/config"
	"github.com/pavanrkadave/homies/pkg/database"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	log.Println("Connected to database successfully!")

	// Run migrations
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database setup complete!")
}

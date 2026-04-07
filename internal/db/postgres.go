package db

import (
	"database/sql"
	"fmt"
	"go-service/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	cfg := config.LoadDBConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	fmt.Printf("Connecting to database with DSN: %s\n", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	fmt.Printf("Pinging database...\n")
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	fmt.Printf("Database ping successful\n")

	return db, nil
}

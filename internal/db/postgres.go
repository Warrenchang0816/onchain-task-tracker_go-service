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

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func InitTables(db *sql.DB) error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS nft_orders (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			image TEXT NOT NULL DEFAULT '',
			price TEXT NOT NULL DEFAULT '',
			recipient_wallet TEXT NOT NULL DEFAULT '',
			creator_wallet TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}
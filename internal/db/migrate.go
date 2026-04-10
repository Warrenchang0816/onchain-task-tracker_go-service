package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(db *sql.DB, dir string) error {
	if err := ensureTable(db); err != nil {
		return err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var names []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".sql" {
			names = append(names, f.Name())
		}
	}

	sort.Strings(names)

	for _, name := range names {
		applied, err := isApplied(db, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		sqlBytes, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration failed %s: %w", name, err)
		}

		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			name,
		); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		fmt.Println("applied:", name)
	}

	return nil
}

func ensureTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id SERIAL PRIMARY KEY,
		version TEXT UNIQUE,
		applied_at TIMESTAMPTZ DEFAULT NOW()
	)
	`)
	return err
}

func isApplied(db *sql.DB, version string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=$1)",
		version,
	).Scan(&exists)
	return exists, err
}
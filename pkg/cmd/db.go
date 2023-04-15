package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/jmoiron/sqlx"
)

func completeMigrateOptions(dbPath string) error {
	_, err := os.Stat(dbPath)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat database file: %w", err)
	}

	dir := filepath.Dir(dbPath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database file: %w", err)
	}
	file.Close()

	return nil
}

func validateMigrateOptions(migrationsPath string) error {
	info, err := os.Stat(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to stat database file: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("migrations path is not a directory: %s", migrationsPath)
	}

	// Read the contents of the migrations directory.
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	// Check if the files within the migrations directory end with .up.sql or .down.sql.
	for _, file := range files {
		if file.IsDir() {
			return fmt.Errorf("migrations directory must not contain subdirectories, found: %s", file.Name())
		}

		name := file.Name()
		if strings.HasSuffix(name, ".up.sql") || strings.HasSuffix(name, ".down.sql") {
			continue
		}

		return errors.New("migration files must have .up.sql or .down.sql suffix, found: " + name)
	}

	return nil
}

func RunMigrate(dbPath, migrationsPath string) error {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

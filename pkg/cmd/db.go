package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ibihim/banking-csv-cli/pkg/sql"
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
	files, err := os.ReadDir(migrationsPath)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	db := sql.NewDatabase(&sql.DatabaseOptions{
		URL: defaultDBPath,
	})
	if err := db.Connect(ctx); err != nil {
		return fmt.Errorf("failed on db connect: %w", err)
	}
	defer db.Close()

	driver, err := db.Driver()
	if err != nil {
		return fmt.Errorf("failed to get db driver: %w", err)
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

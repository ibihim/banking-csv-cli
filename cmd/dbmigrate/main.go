package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	m, err := migrate.New(
		"file://pkg/sql/migrations",
		"sqlite3://transactions.db",
	)
	if err != nil {
		panic(err)
	}

	if err := m.Steps(1); err != nil {
		panic(err)
	}
}

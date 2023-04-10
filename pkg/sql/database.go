package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db                    *sql.DB
	url                   string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
}

type DatabaseOptions struct {
	URL                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

func NewDatabase(opts DatabaseOptions) *Database {
	// - Set WAL mode (not strictly necessary each time because it's persisted in the database, but good for first run)
	// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
	// - Enable foreign key checks
	if opts.URL == "" {
		opts.URL = "./transactions.db"
	}

	opts.URL += "?_journal=WAL&_timeout=5000&_fk=true"

	return &Database{
		url:                   opts.URL,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
	}
}

func (d *Database) Connect(ctx context.Context) error {
	db, err := sql.Open("sqlite3", d.url)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	db.SetMaxOpenConns(d.maxOpenConnections)
	db.SetMaxIdleConns(d.maxIdleConnections)
	db.SetConnMaxLifetime(d.connectionMaxLifetime)
	db.SetConnMaxIdleTime(d.connectionMaxIdleTime)

	d.db = db

	return db.PingContext(ctx)
}

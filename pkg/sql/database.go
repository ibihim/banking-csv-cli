package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ibihim/banking-csv-cli/pkg/transactions"
)

// Database contains a database connection.
type Database struct {
	db                    *sql.DB
	url                   string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
}

// DatabaseOptions contains options for the database.
type DatabaseOptions struct {
	URL                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

// NewDatabase creates a new database.
func NewDatabase(opts *DatabaseOptions) *Database {
	// Set default options
	if opts == nil {
		opts = &DatabaseOptions{}
	}

	// - Set WAL mode (not strictly necessary each time because it's persisted in the database, but good for first run)
	// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
	// - Enable foreign key checks
	if opts.URL == "" {
		opts.URL = "./transactions.db"
	}
	opts.URL += "?cache=shared&mode=rwc&_journal=WAL&_timeout=5000&_fk=true"

	if opts.MaxOpenConnections == 0 {
		opts.MaxOpenConnections = 1
	}

	if opts.MaxIdleConnections == 0 {
		opts.MaxIdleConnections = 1
	}

	if opts.ConnectionMaxLifetime == 0 {
		opts.ConnectionMaxLifetime = -1
	}

	if opts.ConnectionMaxIdleTime == 0 {
		opts.ConnectionMaxIdleTime = -1
	}

	return &Database{
		url:                   opts.URL,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
	}
}

// Connect connects to the database.
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

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

// AddTransaction adds a transaction to the database
func (d *Database) AddTransaction(t *transactions.Transaction) (int64, error) {
	query := `
		INSERT INTO transactions (
			account, booking_date, valuta_date, booking_text, purpose, creditor_id,
			mandate_ref, customer_ref, collector_ref, orig_amount, chargeback_fee,
			beneficiary, account_number, bic, amount, currency, additional_details
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := d.db.Exec(query,
		t.Account, t.BookingDate, t.ValutaDate, t.BookingText, t.Purpose, t.CreditorID,
		t.MandateRef, t.CustomerRef, t.CollectorRef, t.OrigAmount, t.ChargebackFee,
		t.Beneficiary, t.AccountNumber, t.BIC, t.Amount, t.Currency, t.AdditionalDetails,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetTransactions retrieves all transactions from the database
func (d *Database) GetTransactions() ([]*transactions.Transaction, error) {
	query := "SELECT * FROM transactions"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []*transactions.Transaction
	for rows.Next() {
		t := transactions.Transaction{}
		err := rows.Scan(
			&t.ID, &t.Account, &t.BookingDate, &t.ValutaDate, &t.BookingText, &t.Purpose, &t.CreditorID,
			&t.MandateRef, &t.CustomerRef, &t.CollectorRef, &t.OrigAmount, &t.ChargebackFee,
			&t.Beneficiary, &t.AccountNumber, &t.BIC, &t.Amount, &t.Currency, &t.AdditionalDetails,
		)
		if err != nil {
			return nil, err
		}
		ts = append(ts, &t)
	}

	return ts, nil
}

// HasTransaction checks if a transaction already exists in the database
func (d *Database) HasTransaction(transaction *transactions.Transaction) (bool, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE account = ? AND booking_date = ? AND valuta_date = ? AND amount = ? AND creditor_id = ? AND mandate_ref = ?`

	var count int

	if err := d.db.QueryRow(query,
		transaction.Account,
		transaction.BookingDate,
		transaction.ValutaDate,
		transaction.Amount,
		transaction.CreditorID,
		transaction.MandateRef,
	).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

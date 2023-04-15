package cmd

import "github.com/ibihim/banking-csv-cli/pkg/transactions"

type TransactionDatastore interface {
	AddTransaction(transaction *transactions.Transaction) (int64, error)
	GetTransactions() ([]transactions.Transaction, error)
	TransactionExists(transaction *transactions.Transaction) (bool, error)
}

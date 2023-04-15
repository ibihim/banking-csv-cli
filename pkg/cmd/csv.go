package cmd

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/ibihim/banking-csv-cli/pkg/transactions"
)

// ParseTransactions parses a CSV file in CAMT format and returns a slice of Transaction objects
func ParseTransactions(reader io.Reader) ([]*transactions.Transaction, error) {
	// Create a new CSV reader
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	// Read the CSV records
	var ts []*transactions.Transaction
	_, err := csvReader.Read() // skip the header row
	if err != nil {
		return nil, err
	}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Parse the transaction
		transaction := &transactions.Transaction{
			Account:           record[0],
			BookingDate:       record[1],
			ValutaDate:        record[2],
			BookingText:       record[3],
			Purpose:           record[4],
			CreditorID:        record[5],
			MandateRef:        record[6],
			CustomerRef:       record[7],
			CollectorRef:      record[8],
			Beneficiary:       record[11],
			AccountNumber:     record[12],
			BIC:               record[13],
			Currency:          record[15],
			AdditionalDetails: record[16],
		}

		transaction.OrigAmount, err = parseFloat(record[9])
		if err != nil {
			return nil, err
		}
		transaction.ChargebackFee, err = parseFloat(record[10])
		if err != nil {
			return nil, err
		}
		transaction.Amount, err = parseFloat(record[14])
		if err != nil {
			return nil, err
		}

		// Add the transaction to the list
		ts = append(ts, transaction)
	}

	return ts, nil
}

func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
}

// GroupTransactions groups a slice of Transaction objects by beneficiary and calculates the total amount for each beneficiary
func GroupTransactions(transactions []*transactions.Transaction) map[string]float64 {
	groups := make(map[string]float64)
	for _, t := range transactions {
		if t.Amount < 0 {
			groups[t.Beneficiary] -= t.Amount
		} else {
			groups[t.Beneficiary] += t.Amount
		}
	}
	return groups
}

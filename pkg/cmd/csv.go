package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/ibihim/banking-csv-cli/pkg/transactions"
	"k8s.io/klog"
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

		bookingDate, err := time.Parse("02.01.06", record[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse booking date (%q, in record %q): %w", record[1], record, err)
		}

		valutaDate, err := time.Parse("02.01.06", record[2])
		if err != nil {
			valutaDate = bookingDate
			klog.Warningf("failed to parse valuta date (%q, in record %q): %w", record[2], record, err)
		}

		// Parse the transaction
		transaction := &transactions.Transaction{
			Account:           record[0],
			BookingDate:       bookingDate,
			ValutaDate:        valutaDate,
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

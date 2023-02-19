package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Transaction struct {
	Account           string
	BookingDate       string
	ValutaDate        string
	BookingText       string
	Purpose           string
	CreditorID        string
	MandateRef        string
	CustomerRef       string
	CollectorRef      string
	OrigAmount        float64
	ChargebackFee     float64
	Beneficiary       string
	AccountNumber     string
	BIC               string
	Amount            float64
	Currency          string
	AdditionalDetails string
}

type Summary struct {
	Beneficiary string
	TotalAmount float64
}

func main() {
	// Parse the command line arguments
	filename := parseCommandLine()

	// Get the CSV reader
	reader, err := GetReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// Parse the transactions
	transactions, err := ParseTransactions(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Group the transactions by beneficiary
	summaries := GroupTransactions(transactions)

	// Visualize the summaries
	VisualizeTransactions(summaries)
}

func parseCommandLine() string {
	// Define command line flags
	filenamePtr := flag.String("f", "", "the CSV file to parse")
	filePtr := flag.String("file", "", "the CSV file to parse")

	// Parse the command line arguments
	flag.Parse()

	// Check that a file was specified
	if *filenamePtr != "" && *filePtr != "" {
		log.Fatal("Error: both -f and --file flags are specified")
	}
	if *filenamePtr == "" && *filePtr == "" {
		log.Fatal("Error: no input file specified")
	}

	if *filenamePtr != "" {
		return *filenamePtr
	}

	return *filePtr
}

// GetReader returns an io.ReadCloser for the specified file
func GetReader(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

// ParseTransactions parses a CSV file in CAMT format and returns a slice of Transaction objects
func ParseTransactions(reader io.Reader) ([]Transaction, error) {
	// Create a new CSV reader
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	// Read the CSV records
	var transactions []Transaction
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
		transaction := Transaction{
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
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// GroupTransactions groups a slice of Transaction objects by beneficiary and calculates the total amount for each beneficiary
func GroupTransactions(transactions []Transaction) []Summary {
	// Group the transactions by beneficiary
	summaries := make(map[string]Summary)
	for _, transaction := range transactions {
		summary, ok := summaries[transaction.Beneficiary]
		if !ok {
			summary = Summary{Beneficiary: transaction.Beneficiary}
		}
		summary.TotalAmount += transaction.Amount
		summaries[transaction.Beneficiary] = summary
	}

	// Convert the map to a slice
	var result []Summary
	for _, summary := range summaries {
		result = append(result, summary)
	}

	return result
}

// VisualizeTransactions prints a list of Summary objects to the console in a tabular format
func VisualizeTransactions(summaries []Summary) {
	// Print the summaries
	fmt.Printf("%-50s %-10s\n", "Beneficiary", "Total Amount")
	for _, summary := range summaries {
		fmt.Printf("%-50s %-10.2f\n", summary.Beneficiary, summary.TotalAmount)
	}
}

func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
}

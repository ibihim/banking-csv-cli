package cmd

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/ibihim/banking-csv-cli/pkg/model"
	"github.com/ibihim/banking-csv-cli/pkg/transactions"
)

const (
	filename = "filename"
)

func BankingCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "banking",
		Short: "A tool to parse banking csv files",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.CommandLine.VisitAll(func(flag *flag.Flag) {
				klog.V(4).Infof("Flag: --%s=%q", flag.Name, flag.Value)
			})
		},
	}

	// Init klog files
	fs := flag.NewFlagSet("", flag.PanicOnError)
	klog.InitFlags(fs)
	rootCmd.PersistentFlags().AddGoFlagSet(fs)

	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Group transactions by purpose",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename, err := cmd.Flags().GetString(filename)
			if err != nil {
				return err
			}

			return RunGroup(filename)
		},
	}

	groupCmd.Flags().String(filename, "Berlin", "City name")

	rootCmd.AddCommand(groupCmd)

	return rootCmd
}

func RunGroup(filename string) error {
	reader, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer reader.Close()

	transactions, err := ParseTransactions(reader)
	if err != nil {
		return fmt.Errorf("failed to parse transactions: %w", err)
	}

	// Group the transactions by beneficiary and visualize them using Bubble Tea
	groups := GroupTransactions(transactions)

	// Create the table Model
	tableModel := model.MapGroupsToModel(groups)

	if _, err := tea.NewProgram(tableModel).Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}

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

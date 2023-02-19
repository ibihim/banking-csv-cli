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

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Transaction struct {
	// Account is the account number of the account under view.
	Account string
	// BookingDate is the date of the transaction being triggered.
	BookingDate string
	// ValutaDate is the date of the transaction being completed.
	ValutaDate string
	// BookingText is a text that tries to set a type for the transaction.
	BookingText string
	// Purpose is a text that describes the purpose of the transaction.
	Purpose string
	// CreditorID is the creditor identifier of the creditor.
	CreditorID string
	// MandateRef is the id that identifies the mandate that allows the creditor to collect the amount.
	MandateRef string
	// CustomerRef ???
	CustomerRef string
	// CollectorRef is some id that seems to be specific to the creditor.
	CollectorRef string
	// OrigAmount ???
	OrigAmount float64
	// ChargebackFee is the fee charged by the bank for a chargeback.
	ChargebackFee float64
	// Beneficiary is the name of the creditor.
	Beneficiary string
	// AccountNumber is the IBAN of the beneficiary.
	AccountNumber string
	// BIC is the Bank Identifier Code of the beneficiary.
	BIC string
	// Amount is the amount of the transaction.
	Amount float64
	// Currency is the currency of the transaction.
	Currency string
	// AdditionalDetails describes the current state of the transaction.
	AdditionalDetails string
}

type Model struct {
	Groups map[string]float64

	table table.Model
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

	// Group the transactions by beneficiary and visualize them using Bubble Tea
	groups := GroupTransactions(transactions)

	// Create the table Model
	tableModel := mapGroupsToModel(groups)

	if _, err := tea.NewProgram(tableModel).Run(); err != nil {
		log.Fatal(err)
	}
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

func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
}

// GroupTransactions groups a slice of Transaction objects by beneficiary and calculates the total amount for each beneficiary
func GroupTransactions(transactions []Transaction) map[string]float64 {
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

func mapGroupsToModel(groups map[string]float64) Model {
	columns := []table.Column{
		{
			Title: "Beneficiary",
			Width: 50,
		},
		{
			Title: "Amount",
			Width: 10,
		},
	}

	rows := []table.Row{}
	for beneficiary, amount := range groups {
		rows = append(rows, table.Row{beneficiary, fmt.Sprintf("%.2f", amount)})
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithStyles(s),
	)

	return Model{Groups: groups, table: t}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

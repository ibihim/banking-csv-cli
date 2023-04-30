package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ibihim/banking-csv-cli/pkg/table"
	"github.com/ibihim/banking-csv-cli/pkg/transactions"
)

func RunApp(ts []*transactions.Transaction) error {
	summary := transactions.NewSummary(ts)
	table := table.NewTable(summary)

	if _, err := tea.NewProgram(table).Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}

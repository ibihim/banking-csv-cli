package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ibihim/banking-csv-cli/pkg/transactions"
)

func RunApp(ts []*transactions.Transaction) error {
	// Create the table Model
	tableModel := transactions.NewModel(ts)

	if _, err := tea.NewProgram(tableModel).Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}

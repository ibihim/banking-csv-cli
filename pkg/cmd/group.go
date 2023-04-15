package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ibihim/banking-csv-cli/pkg/model"
	"github.com/ibihim/banking-csv-cli/pkg/transactions"
)

func RunGroup(transactions []*transactions.Transaction) error {
	// Group the transactions by beneficiary and visualize them using Bubble Tea
	groups := GroupTransactions(transactions)

	// Create the table Model
	tableModel := model.MapGroupsToModel(groups)

	if _, err := tea.NewProgram(tableModel).Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}

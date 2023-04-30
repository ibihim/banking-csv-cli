package table

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/ibihim/banking-csv-cli/pkg/transactions"
)

type Table struct {
	// rows is a list of rows in the table used to show the bubbles.table.
	rows []table.Row

	// ref is a reference from row to an individual sum witin the model.
	ref []*transactions.Sum

	// model is the model that is used to build the table it represents the
	// state of the application.
	model *transactions.Sum

	// table is the table that is used to display the data.
	table table.Model
}

func NewTable(summary *transactions.Sum) *Table {
	t := &Table{
		rows:  []table.Row{},
		ref:   []*transactions.Sum{},
		model: summary,
	}

	t.buildTable()

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

	t.table = table.New(
		table.WithColumns(createColumns()),
		table.WithRows(t.rows),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithStyles(s),
	)

	return t
}

func (t *Table) action(row int, action string) error {
	//	fmt.Printf("action: %s, row: %d, title: %s", action, row, t.ref[row].Title())

	if err := t.ref[row].Action(action); err != nil {
		return err
	}

	t.buildTable()

	return nil
}

func newRow(date, beneficiary, description string, sum float64) table.Row {
	return table.Row([]string{
		date,
		beneficiary,
		description,
		strconv.FormatFloat(sum, 'f', 2, 64),
	})
}

func (t *Table) reset() {
	t.rows = []table.Row{}
	t.ref = []*transactions.Sum{}
}

func (t *Table) buildTable() *Table {
	t.reset()

	for _, year := range t.model.Sums() {
		if !year.Visible() {
			continue
		}

		t.ref = append(t.ref, year)
		t.rows = append(t.rows, newRow(year.Title(), "", "", year.Total()))

		for _, month := range year.Sums() {
			if !month.Visible() {
				continue
			}

			t.ref = append(t.ref, month)
			t.rows = append(t.rows, newRow(fmt.Sprintf("- %s", month.Title()), "", "", month.Total()))

			for _, beneficiary := range month.Sums() {
				if !beneficiary.Visible() {
					continue
				}

				t.ref = append(t.ref, beneficiary)
				t.rows = append(t.rows, newRow("", beneficiary.Title(), "", beneficiary.Total()))

				for _, transaction := range beneficiary.Sums() {
					if !transaction.Visible() {
						continue
					}

					t.ref = append(t.ref, transaction)
					t.rows = append(t.rows, newRow("", "", transaction.Title(), transaction.Total()))
				}
			}
		}
	}

	return t
}

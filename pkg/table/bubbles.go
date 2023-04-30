package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

const (
	columnCount = 4

	columnWidth1 = 10
	columnWidth2 = 50
	columnWidth3 = 50
	columnWidth4 = 10
)

func createColumns() []table.Column {
	return []table.Column{
		{Title: "Date", Width: columnWidth1},
		{Title: "Beneficiary", Width: columnWidth2},
		{Title: "Description", Width: columnWidth3},
		{Title: "Sum", Width: columnWidth4},
	}
}

func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if t.table.Focused() {
				t.table.Blur()
			} else {
				t.table.Focus()
			}

		case "q", "ctrl+c":
			return t, tea.Quit

		case "enter":
			currentRow := t.table.Cursor()
			t.action(currentRow, "enter")
			t.table.SetRows(t.rows)
		}
	}

	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

func (t *Table) View() string {
	return t.table.View()
}

func (t *Table) Init() tea.Cmd {
	return nil
}

package transactions

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	table table.Model

	window *Window

	rows    []table.Row
	actions []func()
}

func (m Model) updateViewState() {
	clickRows := m.window.GetRows()
	m.rows = make([]table.Row, len(clickRows))
	m.actions = make([]func(), len(clickRows))

	for i, row := range clickRows {
		m.rows[i] = row.GetRow()
		m.actions[i] = row.OnClick
	}
}

func NewModel(transactions []*Transaction) *Model {
	m := Model{
		window: NewWindow(transactions),
	}
	m.updateViewState()

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

	m.table = table.New(
		table.WithColumns(CreateColumns()),
		table.WithRows(m.rows),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithStyles(s),
	)

	return &m
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
			currentRow := m.table.Cursor()
			m.actions[currentRow]()
			m.updateViewState()
			m.table.SetRows(m.rows)
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

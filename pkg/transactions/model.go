package transactions

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	window *Window

	table table.Model
}

func NewModel(transactions []*Transaction) *Model {
	m := Model{
		window: NewWindow(transactions),
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

	m.table = table.New(
		table.WithColumns(CreateColumns()),
		table.WithRows(m.window.GetRows()),
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
			cmds := []tea.Cmd{tea.Printf("TableRow (%d): ", m.table.Cursor())}

			for i, r := range m.table.SelectedRow() {
				cmds = append(cmds, tea.Printf("%d: %s ", i, r))
			}
			cmds = append(cmds, tea.Printf("\n"))

			return m, tea.Sequence(cmds...)
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

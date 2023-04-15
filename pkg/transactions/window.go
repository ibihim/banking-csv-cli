package transactions

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
)

// TODO@ibihim: redo as tree?
// TODO@ibihim: streamling export / unexport

const (
	columnCount = 4

	columnWidth1 = 10
	columnWidth2 = 50
	columnWidth3 = 50
	columnWidth4 = 10
)

func CreateColumns() []table.Column {
	return []table.Column{
		{Title: "Date", Width: columnWidth1},
		{Title: "Beneficiary", Width: columnWidth2},
		{Title: "Description", Width: columnWidth3},
		{Title: "Sum", Width: columnWidth4},
	}
}

func CreateRow(date, beneficiary, description, sum string) table.Row {
	return table.Row{date, beneficiary, description, sum}
}

type Window struct {
	title        string
	transactions []*Transaction

	orderedTabs []*Tab
	tabs        map[string]*Tab
}

func NewWindow(transactions []*Transaction) *Window {
	w := Window{
		title:        "Transactions",
		transactions: transactions,

		orderedTabs: []*Tab{},
		tabs:        make(map[string]*Tab),
	}

	for _, t := range transactions {
		yearStr := strconv.Itoa(t.ValutaDate.Year())
		if !w.HasTab(yearStr) {
			w.AddTab(NewTab(yearStr))
		}
		yearTab := w.GetTab(yearStr)

		monthStr := t.ValutaDate.Month().String()
		if !yearTab.HasCategory(monthStr) {
			yearTab.AddCategory(NewCategory(monthStr))
		}
		monthCategory := yearTab.GetCategory(monthStr)

		if !monthCategory.HasSummary(t.Beneficiary) {
			monthCategory.AddSummary(NewSummary(t.Beneficiary))
		}
		beneficiarySummary := monthCategory.GetSummary(t.Beneficiary)
		beneficiarySummary.AddTransaction(t)
	}

	return &w
}

func (w *Window) HasTab(title string) bool {
	_, ok := w.tabs[title]
	return ok
}

func (w *Window) AddTab(t *Tab) {
	w.tabs[t.title] = t
	w.orderedTabs = append(w.orderedTabs, t)
}

func (w *Window) GetTab(title string) *Tab {
	return w.tabs[title]
}

func (w *Window) GetRows() []table.Row {
	rows := []table.Row{}

	for _, t := range w.orderedTabs {
		rows = append(rows, t.GetRows()...)
	}

	return rows
}

type Tab struct {
	title             string
	orderedCategories []*Category
	categories        map[string]*Category
}

func NewTab(title string) *Tab {
	return &Tab{
		title:             title,
		orderedCategories: []*Category{},
		categories:        make(map[string]*Category),
	}
}

func (t *Tab) HasCategory(title string) bool {
	for _, c := range t.categories {
		if c.title == title {
			return true
		}
	}
	return false
}

func (t *Tab) AddCategory(c *Category) {
	t.categories[c.title] = c
	t.orderedCategories = append(t.orderedCategories, c)
}

func (t *Tab) GetCategory(title string) *Category {
	return t.categories[title]
}

func (t *Tab) GetRows() []table.Row {
	rows := []table.Row{
		CreateRow(fmt.Sprintf("[ %s ]", t.title), "", "", ""),
	}

	for _, c := range t.orderedCategories {
		rows = append(rows, c.GetRows()...)
	}
	return rows
}

type Category struct {
	title            string
	sum              float64
	orderedSummaries []*Summary
	summaries        map[string]*Summary
}

func NewCategory(title string) *Category {
	return &Category{
		title:            title,
		orderedSummaries: []*Summary{},
		summaries:        make(map[string]*Summary),
	}
}

func (c *Category) HasSummary(title string) bool {
	_, ok := c.summaries[title]
	return ok
}

func (c *Category) AddSummary(s *Summary) {
	c.summaries[s.title] = s
	c.orderedSummaries = append(c.orderedSummaries, s)
}

func (c *Category) GetSummary(title string) *Summary {
	return c.summaries[title]
}

func (c *Category) GetSum() float64 {
	c.sum = 0

	for _, s := range c.orderedSummaries {
		c.sum += s.GetSum()
	}

	return c.sum
}

func (c *Category) GetRows() []table.Row {
	rows := []table.Row{CreateRow(
		fmt.Sprintf("- %s", c.title),
		"",
		"",
		fmt.Sprintf("%.2f", c.GetSum()),
	)}

	for _, s := range c.orderedSummaries {
		rows = append(rows, s.GetRows()...)
	}

	return rows
}

type Summary struct {
	title       string
	sum         float64
	showDetails bool

	transactions []*Transaction
}

func NewSummary(title string) *Summary {
	return &Summary{
		title:        title,
		transactions: []*Transaction{},
	}
}

func (s *Summary) AddTransaction(t *Transaction) {
	s.transactions = append(s.transactions, t)
}

func (s *Summary) GetSum() float64 {
	s.sum = 0

	for _, t := range s.transactions {
		s.sum += t.Amount
	}

	return s.sum
}

func (s *Summary) OnClick() {
	s.showDetails = !s.showDetails
}

func (s *Summary) GetRows() []table.Row {
	rows := []table.Row{CreateRow(
		"",
		s.title,
		"",
		fmt.Sprintf("%.2f", s.GetSum()),
	)}

	if !s.showDetails {
		return rows
	}

	// TODO@ibihim sort transactions

	for _, t := range s.transactions {
		rows = append(rows, CreateRow(
			"",
			t.ValutaDate.Format("- 02.01"),
			t.Purpose,
			fmt.Sprintf("%.2f", t.Amount),
		))
	}

	return rows
}

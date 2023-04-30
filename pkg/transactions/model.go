package transactions

import (
	"strconv"
)

func noop(*Sum, string) error { return nil }

func NewSummary(ts []*Transaction) *Sum {
	sum := &Sum{
		title: "Transactions",

		visible: true,
		action:  noop,

		orderedSums: []*Sum{},
		mappedSums:  map[string]*Sum{},
	}

	for _, t := range ts {
		yearStr := strconv.Itoa(t.ValutaDate.Year())
		if !sum.Has(yearStr) {
			sum.AddSum(NewSum(yearStr))
		}
		year := sum.Sum(yearStr)

		monthStr := t.ValutaDate.Month().String()
		if !year.Has(monthStr) {
			year.AddSum(NewSum(monthStr))
		}
		month := year.Sum(monthStr)

		if !month.Has(t.Beneficiary) {
			month.AddSum(NewSum(t.Beneficiary))
		}
		beneficiary := month.Sum(t.Beneficiary)
		beneficiary.AddSum(&Sum{
			title: t.Purpose,
			sum:   t.Amount,

			visible: false,
			action: func(self *Sum, action string) error {
				// Toggle the visibility of the children.
				for _, sum := range self.orderedSums {
					sum.visible = !sum.visible
				}

				return nil
			},

			orderedSums: []*Sum{},
			mappedSums:  map[string]*Sum{},
		})
	}

	return sum
}

func NewSum(title string) *Sum {
	return &Sum{
		title:       title,
		orderedSums: []*Sum{},
		mappedSums:  map[string]*Sum{},
		visible:     true,
		action:      noop,
	}
}

type Sum struct {
	title string
	sum   float64

	visible bool
	action  func(*Sum, string) error

	orderedSums []*Sum // Queue
	mappedSums  map[string]*Sum
}

func (s *Sum) Action(action string) error {
	return s.action(s, action)
}

func (s *Sum) Title() string {
	return s.title
}

func (s *Sum) Visible() bool {
	return s.visible
}

func (s *Sum) AddSum(sum *Sum) {
	s.orderedSums = append(s.orderedSums, sum)
	s.mappedSums[sum.title] = sum
}

func (s *Sum) Sum(title string) *Sum {
	return s.mappedSums[title]
}

func (s *Sum) Sums() []*Sum {
	return s.orderedSums
}

func (s *Sum) Has(title string) bool {
	_, ok := s.mappedSums[title]
	return ok
}

func (s *Sum) Total() float64 {
	// if this is a leaf, return the sum.
	if len(s.orderedSums) == 0 {
		return s.sum
	}

	// if this is a branch, return the sum of all children.
	var total float64
	for _, sum := range s.orderedSums {
		total += sum.Total()
	}
	return total
}

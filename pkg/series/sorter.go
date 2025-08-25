package series

import (
	"sort"
	"time"
)

type By func(t1, t2 *time.Time) bool

type seriesSorter struct {
	s  *TimeSeries
	by By
}

// Len implements sort.Interface.
func (s *seriesSorter) Len() int {
	if s == nil {
		return 0
	}

	return len(s.s.Times)
}

// Less implements sort.Interface.
func (s *seriesSorter) Less(i int, j int) bool {
	return s.by(&s.s.Times[i], &s.s.Times[j])
}

// Swap implements sort.Interface.
func (s *seriesSorter) Swap(i int, j int) {
	s.s.Times[i], s.s.Times[j] = s.s.Times[j], s.s.Times[i]
	s.s.Values[i], s.s.Values[j] = s.s.Values[j], s.s.Values[i]
}

func (by By) Sort(s *TimeSeries) {
	x := &seriesSorter{s: s, by: by}
	sort.Sort(x)
}

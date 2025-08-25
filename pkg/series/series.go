package series

import (
	"errors"
	"time"
)

type TimeSeries struct {
	Times  []time.Time
	Values []float64
}

func (s TimeSeries) Validate() error {
	if len(s.Times) != len(s.Values) {
		return errors.New("fields Times and Values must have the same length")
	}

	return nil
}

func (s TimeSeries) IsSorted() bool {
	if len(s.Times) <= 1 {
		return true
	}

	prev := s.Times[0]
	for _, t := range s.Times[1:] {
		if t.Before(prev) {
			return false
		}

		prev = t
	}

	return true
}

func (s *TimeSeries) SortByDate() {
	date := func(t1, t2 *time.Time) bool {
		if t1 == nil {
			return false
		}

		if t2 == nil {
			return true
		}

		return t1.Before(*t2)
	}

	By(date).Sort(s)
}

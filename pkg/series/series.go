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

package series

import (
	"errors"
	"fmt"
)

func (s *TimeSeries) PeriodReturns() (*TimeSeries, error) {
	if s == nil {
		return nil, errors.New("can't calculate return of nil series")
	}

	if len(s.Times) <= 1 {
		return &TimeSeries{}, nil
	}

	times := s.Times[1:]
	values := make([]float64, 0, len(s.Values)-1)
	prev := s.Values[0]
	for idx, v := range s.Values[1:] {
		if prev == 0 {
			return nil, fmt.Errorf("divide by zero error. s.Values[%d] = 0", idx)
		}

		ret := v/prev - 1
		values = append(values, ret)
		prev = v
	}

	return &TimeSeries{Times: times, Values: values}, nil
}

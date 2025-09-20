package analytics

import (
	"errors"
	"fmt"
	"fundcalc/pkg/series"
)

func PeriodReturns(s *series.TimeSeries) (*series.TimeSeries, error) {
	if s == nil {
		return nil, errors.New("can't calculate return of nil series")
	}

	if len(s.Times) <= 1 {
		return &series.TimeSeries{}, nil
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

	return &series.TimeSeries{Times: times, Values: values}, nil
}

func Returns(periods int) Transform {
	return func(s *series.TimeSeries) (*series.TimeSeries, error) {
		if s == nil {
			return nil, errors.New("can't calculate return of nil series")
		}

		if len(s.Times) <= periods {
			return &series.TimeSeries{}, nil
		}

		times := s.Times[periods:]
		values := make([]float64, 0, len(s.Values)-periods)

		for i := 0; i < len(s.Values)-periods; i++ {
			if s.Values[i] == 0 {
				return nil, fmt.Errorf("divide by zero error. s.Values[%d] = 0", i)
			}

			ret := s.Values[i+periods]/s.Values[i] - 1
			values = append(values, ret)
		}

		return &series.TimeSeries{Times: times, Values: values}, nil
	}
}

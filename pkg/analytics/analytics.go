package analytics

import "fundcalc/pkg/series"

type Transform func(*series.TimeSeries) (*series.TimeSeries, error)

func IdentityTransform(s *series.TimeSeries) (*series.TimeSeries, error) {
	return s, nil
}

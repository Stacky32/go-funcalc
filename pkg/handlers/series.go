package handlers

import (
	"fundcalc/pkg/analytics"
	"fundcalc/pkg/portfolio"
	"fundcalc/pkg/reader"
	"fundcalc/pkg/series"
	"fundcalc/pkg/transformer"
)

func PrepareTimeSeries(weightings portfolio.PortfolioWeightings, funds portfolio.FundDataMap, f analytics.Transform) (*series.TimeSeries, error) {
	labels := make([]portfolio.SeriesKey, 0, len(weightings))
	series := make([]*series.TimeSeries, 0, len(weightings))
	for k := range weightings {
		labels = append(labels, k)
		s := reader.GetPriceSeries(k, funds)
		if s != nil {
			series = append(series, s)
		}
	}

	pivot, err := transformer.Pivot("Portfolio", labels, series)
	if err != nil {
		return nil, err
	}

	combined, err := transformer.CreateWeightedSum(pivot, weightings)
	if err != nil {
		return nil, err
	}

	if err = combined.Validate(); err != nil {
		return nil, err
	}

	combined.SortByDate()
	data, err := f(combined)
	if err != nil {
		return nil, err
	}

	return data, nil
}

package transformer

import (
	"errors"
	"fmt"
	"fundcalc/pkg/series"
	"time"
)

type SeriesKey string

type PortfolioWeightings map[SeriesKey]float64

type DataTable struct {
	Headers map[SeriesKey]struct{}
	Data    map[time.Time]map[SeriesKey]float64
}

func Pivot(key string, labels []SeriesKey, series []*series.TimeSeries) (*DataTable, error) {
	if len(labels) != len(series) {
		return nil, errors.New("labels and series must be the same length")
	}

	if len(series) == 0 {
		return &DataTable{}, nil
	}

	headers := make(map[SeriesKey]struct{}, len(labels))
	for _, l := range labels {
		headers[l] = struct{}{}
	}

	data := make(map[time.Time]map[SeriesKey]float64, len(series[0].Times))

	for i1, s := range series {
		label := labels[i1]
		for i2, t := range s.Times {
			row := data[t]
			if row == nil {
				data[t] = map[SeriesKey]float64{}
			}

			data[t][label] = s.Values[i2]
		}
	}

	// Postprocess
	for ts, r := range data {
		// Check for missing values in row
		for _, v := range r {
			if v == 0 {
				delete(data, ts)
				break
			}
		}
	}

	dt := DataTable{
		Headers: headers,
		Data:    data,
	}

	return &dt, nil
}

func CreateWeightedSum(dt *DataTable, weights PortfolioWeightings) (*series.TimeSeries, error) {
	if dt == nil {
		return nil, errors.New("input data table must not be nil")
	}

	// check weightings exist in headers
	for k := range weights {
		if _, ok := dt.Headers[k]; !ok {
			return nil, fmt.Errorf("fund %s is missing from input data", k)
		}
	}

	acc := func(row map[SeriesKey]float64) float64 {
		total := float64(0)
		for k, v := range row {
			total += weights[k] * v
		}

		return total
	}

	times := make([]time.Time, 0, len(dt.Data))
	values := make([]float64, 0, len(dt.Data))

	for ts, x := range dt.Data {
		times = append(times, ts)
		values = append(values, acc(x))
	}

	return &series.TimeSeries{Times: times, Values: values}, nil
}

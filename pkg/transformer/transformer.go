package transformer

import (
	"errors"
	"fmt"
	"fundcalc/pkg/reader"
)

type TimeStamp = reader.TimeStamp
type SeriesKey string

type DataTable struct {
	Key     SeriesKey
	Data    map[TimeStamp](map[SeriesKey]float32)
	Headers map[SeriesKey]struct{}
}

type PortfolioWeightings map[SeriesKey]float32

type SimpleSeries struct {
	Key  SeriesKey
	Data []reader.DataPoint
}

func (s *SimpleSeries) SortByDate() {
	date := func(p1, p2 *reader.DataPoint) bool {
		d1, err := p1.Date.ParseDate()
		if err != nil {
			return false
		}

		d2, err := p2.Date.ParseDate()
		if err != nil {
			return false
		}

		return d1.Before(d2)
	}

	By(date).Sort(s.Data)
}

func Pivot(key string, series []*SimpleSeries) *DataTable {
	res := &DataTable{
		Key:     SeriesKey(key),
		Data:    map[TimeStamp](map[SeriesKey]float32){},
		Headers: map[SeriesKey]struct{}{},
	}

	for _, s := range series {
		res.Headers[s.Key] = struct{}{}
		for _, p := range s.Data {
			row := res.Data[p.Date]
			if row == nil {
				res.Data[p.Date] = map[SeriesKey]float32{}
			}

			res.Data[p.Date][s.Key] = p.AdjustedClose
		}
	}

	// Postprocess
	for ts, r := range res.Data {
		// Check for missing values in row
		for _, v := range r {
			if v == 0 {
				delete(res.Data, ts)
				break
			}
		}
	}

	return res
}

func CreateWeightedSum(data *DataTable, weights PortfolioWeightings) (*SimpleSeries, error) {
	if data == nil {
		return nil, errors.New("input data table must not be nil")
	}

	// check weightings exist in headers
	for k := range weights {
		if _, ok := data.Headers[k]; !ok {
			return nil, fmt.Errorf("fund %s is missing from input data", k)
		}
	}

	acc := getAccumulator(weights)

	res := &SimpleSeries{Key: SeriesKey(data.Key)}
	for ts, x := range data.Data {
		res.Data = append(res.Data, reader.DataPoint{
			Date:          ts,
			AdjustedClose: acc(x),
		})
	}

	return res, nil
}

type WeightedAccumulator func(map[SeriesKey]float32) float32

func getAccumulator(w PortfolioWeightings) WeightedAccumulator {
	return func(row map[SeriesKey]float32) float32 {
		total := float32(0)
		for k, v := range row {
			total += w[k] * v
		}

		return total
	}
}

package analytics_test

import (
	"errors"
	"fundcalc/pkg/analytics"
	"fundcalc/pkg/series"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GetTestSeries(json string) *series.TimeSeries {
	r := strings.NewReader(json)
	dec := series.NewDecoder(r)
	s := series.TimeSeries{}
	if err := dec.DecodeSeries(&s); err != nil {
		panic(err)
	}

	return &s
}

func SeriesFromValues(start time.Time, vals ...float64) *series.TimeSeries {
	t := start
	s := &series.TimeSeries{Times: make([]time.Time, 0, len(vals)), Values: vals}
	for range vals {
		s.Times = append(s.Times, t)
		t = t.AddDate(0, 0, 1)
	}

	return s
}

func TestPeriodReturns_Properties(t *testing.T) {
	testCases := []struct {
		name          string
		series        *series.TimeSeries
		expectedError error
	}{
		{
			name:          "nil series",
			series:        nil,
			expectedError: errors.New("can't calculate return of nil series"),
		},
		{
			name:   "empty",
			series: &series.TimeSeries{},
		},
		{
			name: "one element",
			series: GetTestSeries(`
			{
				"Times": ["2018-07-04T00:00:00Z"],
				"Values": [105.43]
			}`),
		},
		{
			name: "two elements unordered",
			series: GetTestSeries(`
			{
				"Times": ["2018-07-04T00:00:00Z", "2018-07-03T00:00:00Z"],
				"Values": [107.53, 103.7]
			}`),
		},
		{
			name: "two elements unordered",
			series: GetTestSeries(`
			{
				"Times": ["2018-07-02T00:00:00Z", "2018-07-03T00:00:00Z", "2018-07-04T00:00:00Z", "2018-07-05T00:00:00Z"],
				"Values": [100.00, 101.00, 102.01, 101.90]
			}`),
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.name, func(t *testing.T) {
			rets, err := analytics.PeriodReturns(scenario.series)
			if scenario.expectedError != nil {
				assert.Equal(t, scenario.expectedError, err)
				return
			}

			assert.Nil(t, err)
			assert.Nil(t, rets.Validate())

			if len(scenario.series.Times) <= 1 {
				assert.Equal(t, &series.TimeSeries{}, rets)
			} else {
				assert.Equal(t, len(scenario.series.Times)-1, len(rets.Times))
			}
		})
	}
}

func TestReturns(t *testing.T) {
	from, err := time.Parse("2006-01-02T15:04:05Z", "2018-07-05T00:00:00Z") // Monday
	if err != nil {
		panic(err)
	}

	nilErr := errors.New("can't calculate return of nil series")

	testCases := []struct {
		name     string
		periods  int
		s        *series.TimeSeries
		expected *series.TimeSeries
		err      error
	}{
		{
			name:     "nil series, p=1",
			periods:  1,
			s:        nil,
			expected: nil,
			err:      nilErr,
		},
		{
			name:     "nil series, p=2",
			periods:  1,
			s:        nil,
			expected: nil,
			err:      nilErr,
		},
		{
			name:     "nil series, p=4",
			periods:  1,
			s:        nil,
			expected: nil,
			err:      nilErr,
		},
		{
			name:     "empty series, p=1",
			periods:  1,
			s:        &series.TimeSeries{},
			expected: &series.TimeSeries{},
			err:      nil,
		},
		{
			name:     "empty series, p=2",
			periods:  1,
			s:        &series.TimeSeries{},
			expected: &series.TimeSeries{},
			err:      nil,
		},
		{
			name:     "empty series, p=4",
			periods:  1,
			s:        &series.TimeSeries{},
			expected: &series.TimeSeries{},
			err:      nil,
		},
		{
			name:     "10% daily returns p=1",
			periods:  1,
			s:        SeriesFromValues(from, 1, 1.1, 1.21, 1.331, 1.4641),
			expected: SeriesFromValues(from.AddDate(0, 0, 1), 0.1, 0.1, 0.1, 0.1),
			err:      nil,
		},
		{
			name:     "10% daily returns p=2",
			periods:  2,
			s:        SeriesFromValues(from, 1, 1.1, 1.21, 1.331, 1.4641),
			expected: SeriesFromValues(from.AddDate(0, 0, 2), 0.21, 0.21, 0.21),
			err:      nil,
		},
		{
			name:     "insufficient data points",
			periods:  5,
			s:        SeriesFromValues(from, 1, 1.1, 1.21, 1.331, 1.4641),
			expected: &series.TimeSeries{},
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, err := analytics.Returns(tc.periods)(tc.s)
			if tc.err != nil {
				assert.Equal(t, tc.err, err)
				assert.Nil(t, x)
				return
			}

			assert.Equal(t, tc.expected.Times, x.Times)
			assert.InDeltaSlice(t, tc.expected.Values, x.Values, 1e-9)
		})
	}
}

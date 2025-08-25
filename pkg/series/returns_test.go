package series_test

import (
	"errors"
	"fundcalc/pkg/series"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeriodReturns(t *testing.T) {
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
			series: getTestSeries(`
			{
				"Times": ["2018-07-04T00:00:00Z"],
				"Values": [105.43]
			}`),
		},
		{
			name: "two elements unordered",
			series: getTestSeries(`
			{
				"Times": ["2018-07-04T00:00:00Z", "2018-07-03T00:00:00Z"],
				"Values": [107.53, 103.7]
			}`),
		},
		{
			name: "two elements unordered",
			series: getTestSeries(`
			{
				"Times": ["2018-07-02T00:00:00Z", "2018-07-03T00:00:00Z", "2018-07-04T00:00:00Z", "2018-07-05T00:00:00Z"],
				"Values": [100.00, 101.00, 102.01, 101.90]
			}`),
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.name, func(t *testing.T) {
			rets, err := scenario.series.PeriodReturns()
			if scenario.expectedError != nil {
				assert.Equal(t, scenario.expectedError, err)
				return
			}

			assert.Nil(t, err)

			err = rets.Validate()
			assert.Nil(t, err)

			if len(scenario.series.Times) <= 1 {
				assert.Equal(t, &series.TimeSeries{}, rets)
			} else {
				assert.Equal(t, len(scenario.series.Times)-1, len(rets.Times))
			}
		})
	}
}

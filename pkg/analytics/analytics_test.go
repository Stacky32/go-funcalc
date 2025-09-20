package analytics_test

import (
	"fundcalc/pkg/analytics"
	"fundcalc/pkg/series"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIdentityTransform(t *testing.T) {
	n := time.Now()
	testCases := []struct {
		name string
		s    series.TimeSeries
	}{
		{
			name: "Empty series",
			s:    series.TimeSeries{},
		},
		{
			name: "Single point series",
			s:    series.TimeSeries{Times: []time.Time{n}, Values: []float64{1}},
		},
		{
			name: "Multiple points",
			s: series.TimeSeries{
				Times:  []time.Time{n, n.AddDate(0, 0, 1), n.AddDate(0, 0, 2), n.AddDate(0, 0, 3), n.AddDate(0, 0, 4)},
				Values: []float64{1, 1.2, 1.4, 1.6, 1.8},
			},
		},
		{
			name: "Mising Times",
			s:    series.TimeSeries{Values: []float64{1, 2, 3}},
		},
		{
			name: "Mising Values",
			s:    series.TimeSeries{Times: []time.Time{n, n.AddDate(0, 0, 1)}},
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.name, func(t *testing.T) {
			s := &scenario.s
			x, err := analytics.IdentityTransform(s)

			assert.Nil(t, err)
			assert.Equal(t, s, x)
		})
	}
}

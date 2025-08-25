package series_test

import (
	"encoding/json"
	"errors"
	"fundcalc/pkg/series"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Decoder struct {
	dec *json.Decoder
}

func getTestSeries(json string) *series.TimeSeries {
	r := strings.NewReader(json)
	dec := NewDecoder(r)
	s := series.TimeSeries{}
	if err := dec.DecodeSeries(&s); err != nil {
		panic(err)
	}

	return &s
}

func (d Decoder) DecodeSeries(s *series.TimeSeries) error {
	return d.dec.Decode(s)
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{dec: json.NewDecoder(r)}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name string
		s    *series.TimeSeries
		want error
	}{
		{
			name: "More Values than Times",
			s:    getTestSeries(`{"Times": ["2025-08-17T00:00:00Z"], "Values": [105.33, 103.34]}`),
			want: errors.New("fields Times and Values must have the same length"),
		},
		{
			name: "More Times than Values",
			s:    getTestSeries(`{"Times": ["2025-08-17T00:00:00Z", "2025-08-18T00:00:00Z"], "Values": [105.33]}`),
			want: errors.New("fields Times and Values must have the same length"),
		},
		{
			name: "Same number of Times and Values",
			s:    getTestSeries(`{"Times": ["2025-08-17T00:00:00Z", "2025-08-18T00:00:00Z"], "Values": [105.33, 103.34]}`),
			want: nil,
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.name, func(t *testing.T) {
			err := scenario.s.Validate()
			assert.Equal(t, scenario.want, err)
		})
	}
}

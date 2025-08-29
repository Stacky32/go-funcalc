package charts

import (
	"fundcalc/pkg/series"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type ChartOptions struct {
	Title string
}

type ValueMapper func(float64) float64

func IdentityMapping(x float64) float64 { return x }

func CreatePriceChart(s *series.TimeSeries, op ChartOptions, yMap ValueMapper) *charts.Line {
	line := charts.NewLine()
	if s == nil {
		return line
	}

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{Title: op.Title}),
	)

	x, y := generateAxes(s, yMap)
	line.SetXAxis(x).AddSeries("price", y)

	return line
}

func generateAxes(s *series.TimeSeries, yMap ValueMapper) (x []string, y []opts.LineData) {
	l := len(s.Times)
	x = make([]string, 0, l)
	y = make([]opts.LineData, 0, l)

	for idx, t := range s.Times {
		x = append(x, t.Format("2006-01-02"))

		// Convert from pence to pounds
		y = append(y, opts.LineData{Value: yMap(s.Values[idx])})
	}

	return x, y
}

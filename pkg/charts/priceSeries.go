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

func CreatePriceChart(s *series.TimeSeries, op ChartOptions) *charts.Line {
	line := charts.NewLine()
	if s == nil {
		return line
	}

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{Title: op.Title}),
	)

	x, y := generateAxes(s)
	line.SetXAxis(x).AddSeries("price", y)

	return line
}

func generateAxes(s *series.TimeSeries) (x []string, y []opts.LineData) {
	l := len(s.Times)
	x = make([]string, 0, l)
	y = make([]opts.LineData, 0, l)

	for idx, t := range s.Times {
		x = append(x, t.Format("2006-01-02"))

		// Convert from pence to pounds
		y = append(y, opts.LineData{Value: s.Values[idx] / 100})
	}

	return x, y
}

package charts

import (
	"fmt"
	"fundcalc/pkg/reader"
	"fundcalc/pkg/transformer"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func CreatePriceChart(s *transformer.SimpleSeries) *charts.Line {
	line := charts.NewLine()
	if s == nil {
		return line
	}

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Price series for %s", s.Key),
		}),
	)

	x, y := generateAxes(s.Data)
	line.SetXAxis(x).AddSeries("price", y)

	return line
}

func generateAxes(data []reader.DataPoint) (x []string, y []opts.LineData) {
	l := len(data)
	x = make([]string, 0, l)
	y = make([]opts.LineData, 0, l)

	for _, p := range data {
		x = append(x, string(p.Date))
		// Convert from pence to pounds
		y = append(y, opts.LineData{Value: p.AdjustedClose / 100})
	}

	return x, y
}

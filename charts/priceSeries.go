package charts

import (
	"fmt"
	"fundcalc/reader"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func CreatePriceChart(data []reader.DataPoint, fundName string) *charts.Line {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Price series for %s", fundName),
		}),
	)

	x, y := generateAxes(data)
	line.SetXAxis(x).AddSeries("price", y)

	return line
}

func generateAxes(data []reader.DataPoint) (x []string, y []opts.LineData) {
	l := len(data)
	x = make([]string, 0, l)
	y = make([]opts.LineData, 0, l)

	for _, p := range data {
		x = append(x, p.Date)
		y = append(y, opts.LineData{Value: p.AdjustedClose})
	}

	return x, y
}

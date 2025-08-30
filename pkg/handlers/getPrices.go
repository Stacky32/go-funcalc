package handlers

import (
	"fundcalc/pkg/analytics"
	"fundcalc/pkg/charts"
	"fundcalc/pkg/portfolio"
	"log"
	"net/http"
)

func GetPrices(weightings portfolio.PortfolioWeightings, funds portfolio.FundDataMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := PrepareTimeSeries(weightings, funds, analytics.IdentityTransform)
		if err != nil {
			msg := "Failed to calculate portfolio price series."
			logError(msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		yMap := func(x float64) float64 { return x / 100 }
		line := charts.CreatePriceChart(data, charts.ChartOptions{Title: "Daily prices for portfolio"}, yMap)
		err = line.Render(w)
		if err != nil {
			msg := "Failed to render chart."
			logError(msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

	}
}

func logError(msg string, err error) {
	log.Printf("%s Error:\n%#v\n", msg, err)
}

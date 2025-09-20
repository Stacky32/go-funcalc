package handlers

import (
	"fundcalc/pkg/analytics"
	"fundcalc/pkg/charts"
	"fundcalc/pkg/portfolio"
	"net/http"
)

func GetReturns(weightings portfolio.PortfolioWeightings, funds portfolio.FundDataMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := PrepareTimeSeries(weightings, funds, analytics.Returns(25))
		if err != nil {
			msg := "Unable to caculate period returns."
			logError(msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		line := charts.CreatePriceChart(data, charts.ChartOptions{Title: "Daily returns for portfolio"}, charts.IdentityMapping)
		err = line.Render(w)
		if err != nil {
			msg := "Failed to render chart."
			logError(msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}
}

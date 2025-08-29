package main

import (
	"fmt"
	"fundcalc/pkg/analytics"
	"fundcalc/pkg/charts"
	"fundcalc/pkg/handlers"
	"fundcalc/pkg/reader"
	"fundcalc/pkg/series"
	"fundcalc/pkg/transformer"
	"log"
	"net/http"
	"strings"
)

var funds analytics.FundDataMap
var weightings transformer.PortfolioWeightings

func init() {
	funds = analytics.FundDataMap{
		"rathbone-global":  analytics.FundData{Key: "rathbone-global", Name: "Rathbone Global", Path: "data/rathbone-global.csv"},
		"fssa-asia-focus":  analytics.FundData{Key: "fssa-asia-focus", Name: "FSSA Asia Focus", Path: "data/fssa-asia-focus.csv"},
		"lg-european":      analytics.FundData{Key: "lg-european", Name: "L&G European", Path: "data/lg-european.csv"},
		"lg-international": analytics.FundData{Key: "lg-international", Name: "L&G International", Path: "data/lg-international.csv"},
		"manglg-japan":     analytics.FundData{Key: "manglg-japan", Name: "Man GLG Japan Core Alpha", Path: "data/manglg-japan.csv"},
		"hl-select":        analytics.FundData{Key: "hl-select", Name: "HL Select", Path: "data/hl-select.csv"},
	}

	weightings = transformer.PortfolioWeightings{
		"rathbone-global":  389.39,
		"fssa-asia-focus":  333.208,
		"lg-european":      138.476,
		"lg-international": 307.269,
		"manglg-japan":     293.275,
		"hl-select":        296.755,
	}
}

func main() {
	http.HandleFunc("GET /", handlers.GetIndex)
	http.HandleFunc("GET /portfolio", handleGetPortfolio)
	http.HandleFunc("GET /portfolio-returns", handleGetReturns)

	fmt.Println("Listening on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Critical error. Shutting down...\nError: %#v\n", err)
	}
}

func handleGetPortfolio(w http.ResponseWriter, req *http.Request) {
	data, err := prepData(weightings, analytics.IdentityTransform)
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

func handleGetReturns(w http.ResponseWriter, req *http.Request) {
	data, err := prepData(weightings, analytics.PeriodReturns)
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

func logError(msg string, err error) {
	log.Printf("%s Error:\n%#v\n", msg, err)
}

func getFundData(ref string) analytics.FundData {
	ref = strings.ToLower(ref)
	ref = strings.TrimPrefix(ref, "/")
	return funds[ref]
}

func getFundSeries(ref transformer.SeriesKey) *series.TimeSeries {
	fund := getFundData(string(ref))
	if fund.Key == "" {
		log.Printf("Invalid path requested: %s", ref)
		return nil
	}

	r := reader.CsvPriceReader{Path: fund.Path}
	data, err := r.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV price series data for %s: %v", fund.Path, err)
		return nil
	}

	return data
}

func prepData(weightings transformer.PortfolioWeightings, f analytics.Transform) (*series.TimeSeries, error) {
	labels := make([]transformer.SeriesKey, 0, len(weightings))
	series := make([]*series.TimeSeries, 0, len(weightings))
	for k := range weightings {
		labels = append(labels, k)
		s := getFundSeries(k)
		if s != nil {
			series = append(series, s)
		}
	}

	pivot, err := transformer.Pivot("Portfolio", labels, series)
	if err != nil {
		return nil, err
	}

	combined, err := transformer.CreateWeightedSum(pivot, weightings)
	if err != nil {
		return nil, err
	}

	if err = combined.Validate(); err != nil {
		return nil, err
	}

	combined.SortByDate()
	data, err := f(combined)
	if err != nil {
		return nil, err
	}

	return data, nil
}

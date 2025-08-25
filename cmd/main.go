package main

import (
	"fmt"
	"fundcalc/pkg/charts"
	"fundcalc/pkg/handlers"
	"fundcalc/pkg/reader"
	"fundcalc/pkg/series"
	"fundcalc/pkg/transformer"
	"log"
	"net/http"
	"strings"
)

var funds FundDataMap

type FundDataMap = map[string]FundData

type FundData struct {
	Key  string
	Name string
	Path string
}

func init() {
	funds = FundDataMap{
		"rathbone-global":  FundData{Key: "rathbone-global", Name: "Rathbone Global", Path: "data/rathbone-global.csv"},
		"fssa-asia-focus":  FundData{Key: "fssa-asia-focus", Name: "FSSA Asia Focus", Path: "data/fssa-asia-focus.csv"},
		"lg-european":      FundData{Key: "lg-european", Name: "L&G European", Path: "data/lg-european.csv"},
		"lg-international": FundData{Key: "lg-international", Name: "L&G International", Path: "data/lg-international.csv"},
		"manglg-japan":     FundData{Key: "manglg-japan", Name: "Man GLG Japan Core Alpha", Path: "data/manglg-japan.csv"},
		"hl-select":        FundData{Key: "hl-select", Name: "HL Select", Path: "data/hl-select.csv"},
	}
}

func main() {
	http.HandleFunc("GET /", handlers.GetIndex)
	http.HandleFunc("GET /portfolio", handleGetPortfolio)
	http.HandleFunc("GET /portfolio-returns", handleGetReturns)

	fmt.Println("Listening on http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}

func handleGetPortfolio(w http.ResponseWriter, req *http.Request) {

	data, err := prepData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	yMap := func(x float64) float64 {
		return x
	}

	line := charts.CreatePriceChart(data, charts.ChartOptions{Title: "Daily prices for portfolio"}, yMap)
	err = line.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetReturns(w http.ResponseWriter, req *http.Request) {

	data, err := prepData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rets, err := data.PeriodReturns()
	if err != nil {
		log.Printf("failed to calculate period returns: %#v", err)
		http.Error(w, "Unable to caculate period returns", http.StatusInternalServerError)
		return
	}

	if err = rets.Validate(); err != nil {
		log.Printf("invalid series: %#v", err)
		http.Error(w, "Unable to caculate period returns", http.StatusInternalServerError)
		return
	}

	if !rets.IsSorted() {
		log.Printf("series is not in ascending order")
		http.Error(w, "unable to calculate period returns", http.StatusInternalServerError)
		return
	}

	yMap := func(x float64) float64 {
		return x
	}

	line := charts.CreatePriceChart(rets, charts.ChartOptions{Title: "Daily returns for portfolio"}, yMap)
	err = line.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getFundData(ref string) FundData {
	ref = strings.ToLower(ref)
	ref = strings.TrimPrefix(ref, "/")
	return funds[ref]
}

func getFundSeries(ref transformer.SeriesKey) *series.TimeSeries {
	fund := getFundData(string(ref))
	if fund.Key == "" {
		log.Printf("Invalid path requested: %v", ref)
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

func prepData() (*series.TimeSeries, error) {
	weightings := transformer.PortfolioWeightings{
		"rathbone-global":  389.39,
		"fssa-asia-focus":  333.208,
		"lg-european":      138.476,
		"lg-international": 307.269,
		"manglg-japan":     293.275,
		"hl-select":        296.755,
	}

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

	return combined, nil
}

package main

import (
	"fmt"
	"fundcalc/pkg/charts"
	"fundcalc/pkg/reader"
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
	fmt.Println("Listening on http://localhost:8081")
	http.HandleFunc("/", httpServer)
	http.ListenAndServe(":8081", nil)
}

func httpServer(w http.ResponseWriter, req *http.Request) {

	data, err := prepData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	line := charts.CreatePriceChart(data)
	line.Render(w)
}

func getFundData(ref string) FundData {
	ref = strings.ToLower(ref)
	ref = strings.TrimPrefix(ref, "/")
	return funds[ref]
}

func getFundSeries(ref transformer.SeriesKey) *transformer.SimpleSeries {
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

	return &transformer.SimpleSeries{
		Key:  transformer.SeriesKey(fund.Key),
		Data: data,
	}
}

func prepData() (*transformer.SimpleSeries, error) {
	weightings := transformer.PortfolioWeightings{
		"rathbone-global":  389.39,
		"fssa-asia-focus":  333.208,
		"lg-european":      138.476,
		"lg-international": 307.269,
		"manglg-japan":     293.275,
		"hl-select":        296.755,
	}

	series := []*transformer.SimpleSeries{}
	for k := range weightings {
		s := getFundSeries(k)
		if s != nil {
			series = append(series, s)
		}
	}

	pivot := transformer.Pivot("Portfolio", series)
	combined, err := transformer.CreateWeightedSum(pivot, weightings)
	if err != nil {
		return nil, err
	}

	combined.SortByDate()

	return combined, nil
}

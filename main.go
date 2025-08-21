package main

import (
	"fmt"
	"fundcalc/charts"
	"fundcalc/reader"
	"log"
	"net/http"
	"strings"
)

var funds FundDataMap

type FundDataMap = map[string]FundData

type FundData struct {
	Name string
	Path string
}

func init() {
	funds = FundDataMap{
		"rathbone-global":  FundData{Name: "Rathbone Global", Path: "data/rathbone-global.csv"},
		"fssa-asia-focus":  FundData{Name: "FSSA Asia Focus", Path: "data/fssa-asia-focus.csv"},
		"lg-european":      FundData{Name: "L&G European", Path: "data/lg-european.csv"},
		"lg-international": FundData{Name: "L&G International", Path: "data/lg-international.csv"},
		"manglg-japan":     FundData{Name: "Man GLG Japan Core Alpha", Path: "data/manglg-japan.csv"},
		"hl-select":        FundData{Name: "HL Select", Path: "data/hl-select.csv"},
	}
}

func main() {
	fmt.Println("Listening on http://localhost:8081")
	http.HandleFunc("/", httpServer)
	http.ListenAndServe(":8081", nil)
}

func httpServer(w http.ResponseWriter, req *http.Request) {

	fund := getFundData(req.URL.Path)
	if fund.Path == "" {
		http.Error(w, "Fund not specified in config", http.StatusNotImplemented)
		log.Printf("Invalid path requested: %v", req.URL.Path)
		return
	}

	r := reader.CsvPriceReader{Path: fund.Path}
	data, err := r.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to read CSV price series data for %s: %v", fund.Path, err)
		return
	}

	line := charts.CreatePriceChart(data, fund.Name)
	line.Render(w)
}

func getFundData(ref string) FundData {
	ref = strings.ToLower(ref)
	ref = strings.TrimPrefix(ref, "/")
	return funds[ref]
}

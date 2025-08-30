package main

import (
	"fmt"
	"fundcalc/pkg/handlers"
	"fundcalc/pkg/portfolio"
	"log"
	"net/http"
)

var funds portfolio.FundDataMap
var weightings portfolio.PortfolioWeightings

func init() {
	funds = portfolio.FundDataMap{
		"rathbone-global":  portfolio.FundData{Key: "rathbone-global", Name: "Rathbone Global", Path: "data/rathbone-global.csv"},
		"fssa-asia-focus":  portfolio.FundData{Key: "fssa-asia-focus", Name: "FSSA Asia Focus", Path: "data/fssa-asia-focus.csv"},
		"lg-european":      portfolio.FundData{Key: "lg-european", Name: "L&G European", Path: "data/lg-european.csv"},
		"lg-international": portfolio.FundData{Key: "lg-international", Name: "L&G International", Path: "data/lg-international.csv"},
		"manglg-japan":     portfolio.FundData{Key: "manglg-japan", Name: "Man GLG Japan Core Alpha", Path: "data/manglg-japan.csv"},
		"hl-select":        portfolio.FundData{Key: "hl-select", Name: "HL Select", Path: "data/hl-select.csv"},
	}

	weightings = portfolio.PortfolioWeightings{
		"rathbone-global":  389.39,
		"fssa-asia-focus":  333.208,
		"lg-european":      138.476,
		"lg-international": 307.269,
		"manglg-japan":     293.275,
		"hl-select":        296.755,
	}
}

func main() {
	http.HandleFunc("GET /portfolio/returns", handlers.GetReturns(weightings, funds))
	http.HandleFunc("GET /portfolio/prices", handlers.GetPrices(weightings, funds))
	http.HandleFunc("GET /", handlers.GetIndex)

	fmt.Println("Listening on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Critical error. Shutting down...\nError: %#v\n", err)
	}
}

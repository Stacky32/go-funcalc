package reader

import (
	"fundcalc/pkg/portfolio"
	"fundcalc/pkg/series"
	"log"
	"strings"
)

func getFundData(ref string, funds portfolio.FundDataMap) portfolio.FundData {
	ref = strings.ToLower(ref)
	ref = strings.TrimPrefix(ref, "/")
	return funds[ref]
}

func GetPriceSeries(ref portfolio.SeriesKey, funds portfolio.FundDataMap) *series.TimeSeries {
	fund := getFundData(string(ref), funds)
	if fund.Key == "" {
		log.Printf("Invalid path requested: %s", ref)
		return nil
	}

	r := CsvPriceReader{Path: fund.Path}
	data, err := r.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV price series data for %s: %v", fund.Path, err)
		return nil
	}

	return data
}

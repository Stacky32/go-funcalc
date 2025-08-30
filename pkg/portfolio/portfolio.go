package portfolio

type FundDataMap = map[string]FundData

type FundData struct {
	Key  string
	Name string
	Path string
}

type SeriesKey string
type PortfolioWeightings map[SeriesKey]float64

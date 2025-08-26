package reader

import (
	"encoding/csv"
	"fmt"
	"fundcalc/pkg/series"
	"io"
	"os"
	"slices"
	"strconv"
	"time"
)

type CsvPriceReader struct {
	Path string
}

func (r *CsvPriceReader) ReadAll() (*series.TimeSeries, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	reader := csv.NewReader(file)

	topRow, err := reader.Read()
	if err != nil {
		return nil, err
	}

	dateIdx := slices.IndexFunc(topRow, func(s string) bool { return s == "Date" })
	if dateIdx == -1 {
		err = fmt.Errorf("missing column 'Date': %s", r.Path)
		return nil, err
	}

	adjCloseIdx := slices.IndexFunc(topRow, func(s string) bool { return s == "Adj Close" })
	if dateIdx == -1 {
		err = fmt.Errorf("missing column 'Adj Close': %s", r.Path)
		return nil, err
	}

	times := []time.Time{}
	values := []float64{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		item := record[adjCloseIdx]
		if item == "" || item == "null" {
			continue
		}

		adjClose, err := strconv.ParseFloat(record[adjCloseIdx], 64)
		if err != nil {
			return nil, err
		}

		t, err := series.TimeStamp(record[dateIdx]).Parse()
		if err != nil {
			return nil, err
		}

		times = append(times, t)
		values = append(values, float64(adjClose))
	}

	return &series.TimeSeries{Times: times, Values: values}, nil
}

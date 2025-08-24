package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

type PriceReader interface {
	ReadAll() ([]DataPoint, error)
}

type CsvPriceReader struct {
	Path string
}

func (r *CsvPriceReader) ReadAll() ([]DataPoint, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

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

	data := []DataPoint{}

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

		adjClose, err := strconv.ParseFloat(record[adjCloseIdx], 32)
		if err != nil {
			return nil, err
		}

		d := DataPoint{Date: TimeStamp(record[dateIdx]), AdjustedClose: float32(adjClose)}
		data = append(data, d)
	}

	return data, nil
}

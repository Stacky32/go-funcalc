package reader

import (
	"errors"
	"time"
)

type TimeStamp string

func (t *TimeStamp) ParseDate() (time.Time, error) {
	format := "2006-01-02"
	if t == nil {
		return time.Time{}, errors.New("failed to parse time stamp: nil")
	}

	parsed, err := time.Parse(format, string(*t))
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

type DataPoint struct {
	Date          TimeStamp
	AdjustedClose float32
}

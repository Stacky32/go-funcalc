package series

import (
	"time"
)

var (
	formatString = "2006-01-02"
)

type TimeStamp string

func (t TimeStamp) Parse() (time.Time, error) {
	tt, err := time.Parse(formatString, string(t))
	if err != nil {
		return time.Time{}, err
	}

	return tt, nil
}

func NewTimeStamp(t time.Time) TimeStamp {
	return TimeStamp(t.Format(formatString))
}

func (t TimeStamp) String() string {
	return string(t)
}

package series

import "time"

type TimeSeries struct {
	Times  []time.Time
	Values []float64
}

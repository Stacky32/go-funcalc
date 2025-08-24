package transformer

import (
	"fundcalc/pkg/reader"
	"sort"
)

// DataPoint sorter
type By func(p1, p2 *reader.DataPoint) bool

func (by By) Sort(points []reader.DataPoint) {
	s := &dataSorter{points: points, by: by}
	sort.Sort(s)
}

type dataSorter struct {
	points []reader.DataPoint
	by     By
}

// Len implements sort.Interface.
func (d *dataSorter) Len() int {
	return len(d.points)
}

// Less implements sort.Interface.
func (d *dataSorter) Less(i int, j int) bool {
	return d.by(&d.points[i], &d.points[j])
}

// Swap implements sort.Interface.
func (d *dataSorter) Swap(i int, j int) {
	d.points[i], d.points[j] = d.points[j], d.points[i]
}

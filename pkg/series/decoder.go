package series

import (
	"encoding/json"
	"io"
)

type Decoder struct {
	dec *json.Decoder
}

func (d Decoder) DecodeSeries(s *TimeSeries) error {
	return d.dec.Decode(s)
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{dec: json.NewDecoder(r)}
}

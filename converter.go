package eurex

import (
	"time"

	"github.com/filiptubic/eurex/ecb"
)

var (
	DefaultConverter = ecb.New("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml")
)

type ConverterInterface interface {
	Convert(date time.Time, value float64, from, to string) (converted float64, err error)
}

type Converter struct {
	converter ConverterInterface
}

func New(converter ConverterInterface) *Converter {
	return &Converter{
		converter: converter,
	}
}

func (c *Converter) Convert(date time.Time, value float64, from, to string) (float64, error) {
	return c.converter.Convert(date, value, from, to)
}

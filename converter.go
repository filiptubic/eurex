package eurex

import (
	"time"

	"github.com/filiptubic/eurex/currency"
	"github.com/filiptubic/eurex/ecb"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultLogger        = log.New()
	DefaultClientOptions = ecb.NewECBOptions(3, time.Second*3)
	DefaultClient        = ecb.NewECBClient("https", "www.ecb.europa.eu", DefaultClientOptions, DefaultLogger)
	DefaultConverter     = ecb.New(DefaultClient, true, DefaultLogger)
)

type ConverterInterface interface {
	Convert(date time.Time, value float64, from, to currency.Currency) (converted float64, err error)
}

type Converter struct {
	converter ConverterInterface
}

func New(converter ConverterInterface) *Converter {
	return &Converter{
		converter: converter,
	}
}

func (c *Converter) Convert(date time.Time, value float64, from, to currency.Currency) (float64, error) {
	return c.converter.Convert(date, value, from, to)
}

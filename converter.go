/*
	This package defines Converter API for money conversion from one currency to another.
	Default implementation uses ECB and more details are available inside ecb package.
*/
package eurex

import (
	"time"

	"github.com/filiptubic/eurex/currency"
	"github.com/filiptubic/eurex/ecb"
	log "github.com/sirupsen/logrus"
)

var (
	// Default logger.
	DefaultLogger = log.New()
	// Default option for ECB client.
	DefaultClientOptions = ecb.NewECBOptions(3, time.Second*3)
	// Default ECB client.
	DefaultClient = ecb.NewECBClient("https", "www.ecb.europa.eu", DefaultClientOptions, DefaultLogger)
	// Default ECB converter.
	DefaultConverter = ecb.New(DefaultClient, true, DefaultLogger)
)

// Converter defines converter API.
type Converter interface {
	Convert(date time.Time, value float64, from, to currency.Currency) (converted float64, err error)
}

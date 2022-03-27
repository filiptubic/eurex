/*
	This package holds everything related to ECB converter. Since ECB updates are rare (once per day), converter client uses caching
	in order to boost perfomance and avoid necessary http calls. Keep in mind that ECB has data for the last 90 days and accessing data
	outside of this timeline will result in error.

	NOTE: Russian RUB currency is removed from ECB and the last EUR/RUB update was on 1. March 2022.
*/
package ecb

import (
	"time"

	"github.com/filiptubic/eurex/currency"
	log "github.com/sirupsen/logrus"
)

type currencyMap map[currency.Currency]float64

// Rates is type used for storing ECB rates.
// As underlaying data structure it uses hash map for quick access to specific currency rate for certain date.
// Additionally, it holds information what is to most earliest/oldest date available.
type Rates struct {
	first, last time.Time
	rates       map[time.Time]currencyMap
}

// ECBConverter is ECB implementation of Converter interface. It supports rates caching for better perfomance.
type ECBConverter struct {
	logger *log.Logger
	cache  bool
	cached *Rates
	client ECBClientInterface
}

// New creates ECBConverter object.
func New(client ECBClientInterface, cache bool, logger *log.Logger) *ECBConverter {
	if logger == nil {
		logger = log.New()
	}
	return &ECBConverter{client: client, cache: cache, logger: logger}
}

// newRates makes Rates object from raw ECBResponseData object.
func (c *ECBConverter) newRates(data *ECBResponseData) (*Rates, error) {
	rates := &Rates{
		rates: make(map[time.Time]currencyMap),
	}

	for _, date := range data.Data {
		t, err := date.Date.toTime()
		if err != nil {
			return nil, err
		}

		// set earliest date
		zeroTime := time.Time{}
		if rates.first == zeroTime {
			rates.first = t
		} else if rates.first.After(t) {
			rates.first = t
		}

		// set latest date
		if rates.last == zeroTime {
			rates.last = t
		} else if rates.last.Before(t) {
			rates.last = t
		}

		rates.rates[t] = make(currencyMap)

		for _, rate := range date.Rates {
			currency, ok := currency.Currencies[rate.Currency]
			if !ok {
				return nil, InvalidCurrency{currency: rate.Currency}
			}
			rates.rates[t][currency] = rate.Rate
		}
	}
	return rates, nil
}

// GetRates fetches rates via ECBClient if rate for certain date is not found in cache or caching is disabled.
// When caching is enabled, and cache is present, new data is added to cache only when queried date is not found inside cache.
func (c *ECBConverter) GetRates(date time.Time) (*Rates, error) {
	if c.cache && c.cached != nil && c.cached.rates != nil {
		// if rates are cached for queried date, just return it, don't make http call
		if _, ok := c.cached.rates[date]; ok {
			c.logger.Debugf("using cached rates for date %v", date)
			return c.cached, nil
		}
	}

	// fetch data from ECB
	data, err := c.client.GetRates()
	if err != nil {
		return nil, err
	}
	rates, err := c.newRates(data)
	if err != nil {
		return nil, err
	}
	if c.cache {
		c.cached = rates
	}

	return rates, nil
}

// Convert converts specified value from one currency to another for certian date.
func (c *ECBConverter) Convert(date time.Time, value float64, from, to currency.Currency) (float64, error) {
	if !IsValidCurrency(from, date) {
		return -1, InvalidCurrency{string(from)}
	}

	if !IsValidCurrency(to, date) {
		return -1, InvalidCurrency{string(to)}
	}

	if from == to {
		return value, nil
	}

	rates, err := c.GetRates(date)
	if err != nil {
		return -1, err
	}

	if date.Before(rates.first) || date.After(rates.last) {
		return -1, DateOutOfBound{date, rates.first, rates.last}
	}

	if _, ok := rates.rates[date][from]; from != currency.EUR && !ok {
		return -1, InvalidCurrency{string(from)}
	}

	if _, ok := rates.rates[date][to]; to != currency.EUR && !ok {
		return -1, InvalidCurrency{string(from)}
	}

	if from == currency.EUR {
		return value * rates.rates[date][to], nil
	}

	if to == currency.EUR {
		return value / rates.rates[date][from], nil
	}

	return (value * rates.rates[date][to]) / rates.rates[date][from], nil
}

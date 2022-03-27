package ecb

import (
	"time"

	"github.com/filiptubic/eurex/currency"
	log "github.com/sirupsen/logrus"
)

type currencyMap map[currency.Currency]float64

type Rates struct {
	first, last time.Time
	rates       map[time.Time]currencyMap
}

type ECBConverter struct {
	logger *log.Logger
	cache  bool
	cached *Rates
	client ECBClientInterface
}

func New(client ECBClientInterface, cache bool, logger *log.Logger) *ECBConverter {
	if logger == nil {
		logger = log.New()
	}
	return &ECBConverter{client: client, cache: cache, logger: logger}
}

func (c *ECBConverter) newRates(data *ECBResponseData) (*Rates, error) {
	rates := &Rates{
		rates: make(map[time.Time]currencyMap),
	}

	for _, date := range data.Data {
		t, err := date.Date.toTime()
		if err != nil {
			return nil, err
		}

		zeroTime := time.Time{}
		if rates.first == zeroTime {
			rates.first = t
		} else if rates.first.After(t) {
			rates.first = t
		}

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

func (c *ECBConverter) GetRates(date time.Time) (*Rates, error) {
	if c.cache && c.cached != nil && c.cached.rates != nil {
		if _, ok := c.cached.rates[date]; ok {
			c.logger.Debugf("using cached rates for date %v", date)
			return c.cached, nil
		}
	}

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

func (c *ECBConverter) Convert(date time.Time, value float64, from, to currency.Currency) (float64, error) {
	if !IsValid(from, date) {
		return -1, InvalidCurrency{string(from)}
	}

	if !IsValid(to, date) {
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

	if from == currency.EUR {
		return value * rates.rates[date][to], nil
	}

	if to == currency.EUR {
		return value / rates.rates[date][from], nil
	}

	return (value * rates.rates[date][to]) / rates.rates[date][from], nil
}

package ecb

import (
	"time"

	"github.com/filiptubic/eurex/currency"
)

type currencyMap map[currency.Currency]float64

type ECBConverter struct {
	client ECBClientInterface
}

func New(client ECBClientInterface) *ECBConverter {
	return &ECBConverter{client: client}
}

func (c *ECBConverter) makeRatesMap(data *ECBResponseData) (map[time.Time]currencyMap, error) {
	rates := make(map[time.Time]currencyMap)
	for _, date := range data.Data {
		t, err := date.Date.toTime()
		if err != nil {
			return nil, err
		}

		rates[t] = make(currencyMap)
		for _, rate := range date.Rates {
			currency, ok := currency.Currencies[rate.Currency]
			if !ok {
				return nil, InvalidCurrency{currency: rate.Currency}
			}
			rates[t][currency] = rate.Rate
		}
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

	data, err := c.client.GetRates()
	if err != nil {
		return -1, err
	}

	ratesMap, err := c.makeRatesMap(data)
	if err != nil {
		return -1, err
	}

	if from == currency.EUR {
		return value * ratesMap[date][to], nil
	}

	if to == currency.EUR {
		return value / ratesMap[date][from], nil
	}

	return (value * ratesMap[date][to]) / ratesMap[date][from], nil
}

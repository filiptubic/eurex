package ecb

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"
)

type currencyMap map[string]float64

type ECBConverter struct {
	url string
}

func New(url string) *ECBConverter {
	return &ECBConverter{
		url: url,
	}
}

func (c *ECBConverter) fetch() (*ECBResponseData, error) {
	resp, err := http.Get(c.url)
	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ecbData := ECBResponseData{}
	xml.Unmarshal(respDataBytes, &ecbData)
	return &ecbData, err
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
			rates[t][rate.Currency] = rate.Rate
		}
	}
	return rates, nil
}

func (c *ECBConverter) Convert(date time.Time, value float64, from, to string) (float64, error) {
	if !IsValid(from, date) {
		return -1, InvalidCurrency{from}
	}

	if !IsValid(to, date) {
		return -1, InvalidCurrency{to}
	}

	if from == to {
		return value, nil
	}

	data, err := c.fetch()
	if err != nil {
		return -1, err
	}

	ratesMap, err := c.makeRatesMap(data)
	if err != nil {
		return -1, err
	}

	if from == string(EUR) {
		return value * ratesMap[date][to], nil
	}

	if to == string(EUR) {
		return value / ratesMap[date][to], nil
	}

	return (value * ratesMap[date][to]) / ratesMap[date][from], nil
}
